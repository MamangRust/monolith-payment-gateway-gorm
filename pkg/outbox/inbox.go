package outbox

import (
	"context"
	"errors"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"gorm.io/gorm"
)

var ErrInvalidInboxKey = errors.New("invalid consumer inbox key")

// ConsumerInbox is the durable deduplication contract used by Kafka handlers.
type ConsumerInbox interface {
	Reserve(ctx context.Context, consumerName, eventKey, topic string, partition int32, offset int64) (bool, bool, int64, error)
	MarkProcessed(ctx context.Context, consumerName, eventKey string, reservationVersion int64) error
	Release(ctx context.Context, consumerName, eventKey string, reservationVersion int64, processingErr error) error
}

type reservationRow struct {
	Reserved           bool  `gorm:"column:reserved"`
	Processed          bool  `gorm:"column:processed"`
	ReservationVersion int64 `gorm:"column:reservation_version"`
}

// Reserve claims an event for a consumer. It returns false when the event was
// already processed. An expired processing lease may be reclaimed after a
// consumer crashes. It renders the same lease-fenced upsert as the original
// ReserveConsumerInbox statement through GORM.
func Reserve(ctx context.Context, tx *gorm.DB, consumerName, eventKey, topic string, partition int32, offset int64) (bool, bool, int64, error) {
	if tx == nil || consumerName == "" || eventKey == "" {
		return false, false, 0, ErrInvalidInboxKey
	}
	var row reservationRow
	err := tx.WithContext(ctx).Raw(`
		WITH reserved AS (
			INSERT INTO consumer_inbox (
				consumer_name, event_key, topic, partition_id, message_offset,
				status, attempts, reservation_version, lease_until, last_error, processed_at
			)
			VALUES (?, ?, ?, ?, ?, 'processing', 1, 1,
			        current_timestamp + interval '1 minute', '', NULL)
			ON CONFLICT (consumer_name, event_key) DO UPDATE
			SET status = 'processing',
			    attempts = consumer_inbox.attempts + 1,
			    reservation_version = consumer_inbox.reservation_version + 1,
			    lease_until = current_timestamp + interval '1 minute',
			    last_error = '',
			    topic = EXCLUDED.topic,
			    partition_id = EXCLUDED.partition_id,
			    message_offset = EXCLUDED.message_offset
			WHERE consumer_inbox.status <> 'processed'
			  AND consumer_inbox.lease_until <= current_timestamp
			RETURNING reservation_version
		)
		SELECT
			EXISTS (SELECT 1 FROM reserved) AS reserved,
			EXISTS (
				SELECT 1
				FROM consumer_inbox ci
				WHERE ci.consumer_name = ?
				  AND ci.event_key = ?
				  AND ci.status = 'processed'
			) AS processed,
			COALESCE(
				(SELECT reservation_version FROM reserved),
				(SELECT ci.reservation_version FROM consumer_inbox ci WHERE ci.consumer_name = ? AND ci.event_key = ?)
			) AS reservation_version`,
		consumerName, eventKey, topic, partition, offset,
		consumerName, eventKey, consumerName, eventKey,
	).Scan(&row).Error
	if err != nil {
		return false, false, 0, err
	}
	return row.Reserved, row.Processed, row.ReservationVersion, nil
}

func MarkProcessed(ctx context.Context, tx *gorm.DB, consumerName, eventKey string, reservationVersion int64) error {
	if tx == nil || consumerName == "" || eventKey == "" {
		return ErrInvalidInboxKey
	}
	// Completes only the active reservation.
	return tx.WithContext(ctx).Model(&models.ConsumerInbox{}).
		Where("consumer_name = ? AND event_key = ? AND status = 'processing' AND reservation_version = ?",
			consumerName, eventKey, reservationVersion).
		Updates(map[string]interface{}{
			"status":       "processed",
			"processed_at": gorm.Expr("current_timestamp"),
			"lease_until":  gorm.Expr("current_timestamp"),
			"last_error":   "",
		}).Error
}

func Release(ctx context.Context, tx *gorm.DB, consumerName, eventKey string, reservationVersion int64, processingErr error) error {
	if tx == nil || consumerName == "" || eventKey == "" {
		return ErrInvalidInboxKey
	}
	lastError := "consumer processing failed"
	if processingErr != nil {
		lastError = processingErr.Error()
	}
	// Releases only the active reservation.
	return tx.WithContext(ctx).Model(&models.ConsumerInbox{}).
		Where("consumer_name = ? AND event_key = ? AND status = 'processing' AND reservation_version = ?",
			consumerName, eventKey, reservationVersion).
		Updates(map[string]interface{}{
			"status":      "pending",
			"lease_until": gorm.Expr("current_timestamp"),
			"last_error":  lastError,
		}).Error
}

// PostgresInbox adapts a GORM database to ConsumerInbox. Reservation and
// completion are committed independently because an external side effect
// cannot share a PostgreSQL transaction with the Kafka consumer.
type PostgresInbox struct {
	db *gorm.DB
}

func NewPostgresInbox(db *gorm.DB) (*PostgresInbox, error) {
	if db == nil {
		return nil, errors.New("inbox database is nil")
	}
	return &PostgresInbox{db: db}, nil
}

func (i *PostgresInbox) Reserve(ctx context.Context, consumerName, eventKey, topic string, partition int32, offset int64) (bool, bool, int64, error) {
	var (
		reserved, processed bool
		version             int64
	)
	err := i.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var err error
		reserved, processed, version, err = Reserve(ctx, tx, consumerName, eventKey, topic, partition, offset)
		return err
	})
	if err != nil {
		return false, false, 0, err
	}
	return reserved, processed, version, nil
}

func (i *PostgresInbox) MarkProcessed(ctx context.Context, consumerName, eventKey string, reservationVersion int64) error {
	return i.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return MarkProcessed(ctx, tx, consumerName, eventKey, reservationVersion)
	})
}

func (i *PostgresInbox) Release(ctx context.Context, consumerName, eventKey string, reservationVersion int64, processingErr error) error {
	return i.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return Release(ctx, tx, consumerName, eventKey, reservationVersion, processingErr)
	})
}
