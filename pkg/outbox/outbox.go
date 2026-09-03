package outbox

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrInvalidEvent = errors.New("invalid outbox event")
	ErrNoPublisher  = errors.New("outbox publisher is nil")
)

// Event is a message that must be persisted in the same database transaction
// as the business mutation that produced it. EventKey is an application-level
// idempotency key and is unique in the outbox table.
type Event struct {
	EventKey string
	Topic    string
	Key      string
	Payload  []byte
}

func (e Event) Validate() error {
	if e.EventKey == "" || e.Topic == "" || len(e.Payload) == 0 || !json.Valid(e.Payload) {
		return ErrInvalidEvent
	}
	return nil
}

// Executor is satisfied by pgx.Tx and keeps enqueue usable inside an existing
// business transaction. The caller owns commit/rollback of the transaction.
type Executor interface {
	db.DBTX
}

// Enqueue persists an event in the caller's transaction. A duplicate event key
// is treated as success, which makes retries of a business transaction safe.
func Enqueue(ctx context.Context, tx Executor, event Event) error {
	if err := event.Validate(); err != nil {
		return err
	}
	if tx == nil {
		return errors.New("outbox executor is nil")
	}
	// Enqueue intentionally remains a small generic helper for callers that
	// already own a transaction. Domain commands use SQLC queries for their
	// event insertion; the helper is retained for package-level compatibility.
	return db.New(tx).EnqueueOutboxEvent(ctx, db.EnqueueOutboxEventParams{
		EventKey:   event.EventKey,
		Topic:      event.Topic,
		MessageKey: event.Key,
		Column4:    event.Payload,
	})
}

// Publisher is implemented by pkg/kafka.Kafka and deliberately kept small so
// the outbox package does not depend on the Kafka implementation.
type Publisher interface {
	SendMessageWithRetry(topic, key string, value []byte, attempts int) error
}

type RelayConfig struct {
	BatchSize    int
	PollInterval time.Duration
	MaxAttempts  int
	BaseBackoff  time.Duration
	Topics       []string
}

func (c RelayConfig) withDefaults() RelayConfig {
	if c.BatchSize <= 0 {
		c.BatchSize = 50
	}
	if c.PollInterval <= 0 {
		c.PollInterval = time.Second
	}
	if c.MaxAttempts <= 0 {
		c.MaxAttempts = 10
	}
	if c.BaseBackoff <= 0 {
		c.BaseBackoff = time.Second
	}
	return c
}

type Relay struct {
	pool      *pgxpool.Pool
	publisher Publisher
	config    RelayConfig
}

func NewRelay(pool *pgxpool.Pool, publisher Publisher, config RelayConfig) (*Relay, error) {
	if pool == nil {
		return nil, errors.New("outbox pool is nil")
	}
	if publisher == nil {
		return nil, ErrNoPublisher
	}
	return &Relay{pool: pool, publisher: publisher, config: config.withDefaults()}, nil
}

// Run polls until ctx is cancelled. A claimed event remains recoverable when
// a process crashes: the claim is stored as publishing and stale claims are
// returned to pending by the next claim operation.
func (r *Relay) Run(ctx context.Context) error {
	for {
		processed, err := r.runBatch(ctx)
		if errors.Is(ctx.Err(), context.Canceled) {
			return ctx.Err()
		}
		if err != nil {
			// Keep the relay alive across transient database/Kafka state errors.
			// Claimed rows remain recoverable through the stale-publishing lease.
			if !wait(ctx, r.config.PollInterval) {
				return ctx.Err()
			}
			continue
		}
		if processed == 0 && !wait(ctx, r.config.PollInterval) {
			return ctx.Err()
		}
	}
}

func wait(ctx context.Context, duration time.Duration) bool {
	timer := time.NewTimer(duration)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return false
	case <-timer.C:
		return true
	}
}

type claimedEvent struct {
	ID           int64
	Topic        string
	Key          string
	Payload      []byte
	Attempt      int
	ClaimVersion int64
}

func (r *Relay) runBatch(ctx context.Context) (int, error) {
	events, err := r.claim(ctx)
	if err != nil {
		return 0, err
	}
	for _, event := range events {
		if err := r.publish(ctx, event); err != nil {
			if updateErr := r.markFailed(ctx, event, err); updateErr != nil {
				return len(events), fmt.Errorf("record outbox failure: %w (publish error: %v)", updateErr, err)
			}
			continue
		}
		if err := r.markPublished(ctx, event); err != nil {
			return len(events), err
		}
	}
	return len(events), nil
}

func (r *Relay) claim(ctx context.Context) ([]claimedEvent, error) {
	queries := db.New(r.pool)
	var rows []*db.ClaimOutboxEventsRow
	var err error
	if len(r.config.Topics) == 0 {
		rows, err = queries.ClaimOutboxEvents(ctx, int32(r.config.BatchSize))
	} else {
		topicRows, topicErr := queries.ClaimOutboxEventsByTopics(ctx, db.ClaimOutboxEventsByTopicsParams{
			Limit:   int32(r.config.BatchSize),
			Column2: r.config.Topics,
		})
		err = topicErr
		if topicErr == nil {
			rows = make([]*db.ClaimOutboxEventsRow, 0, len(topicRows))
			for _, row := range topicRows {
				rows = append(rows, &db.ClaimOutboxEventsRow{
					ID: row.ID, Topic: row.Topic, MessageKey: row.MessageKey,
					Payload: row.Payload, Attempts: row.Attempts, ClaimVersion: row.ClaimVersion,
				})
			}
		}
	}
	if err != nil {
		return nil, err
	}
	events := make([]claimedEvent, 0, len(rows))
	for _, row := range rows {
		events = append(events, claimedEvent{
			ID: row.ID, Topic: row.Topic, Key: row.MessageKey,
			Payload: row.Payload, Attempt: int(row.Attempts), ClaimVersion: row.ClaimVersion,
		})
	}
	return events, nil
}

func (r *Relay) publish(ctx context.Context, event claimedEvent) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return r.publisher.SendMessageWithRetry(event.Topic, event.Key, event.Payload, 1)
}

func (r *Relay) markPublished(ctx context.Context, event claimedEvent) error {
	return db.New(r.pool).MarkOutboxPublished(ctx, db.MarkOutboxPublishedParams{
		ID: event.ID, ClaimVersion: event.ClaimVersion,
	})
}

func (r *Relay) markFailed(ctx context.Context, event claimedEvent, publishErr error) error {
	status := "pending"
	if event.Attempt >= r.config.MaxAttempts {
		status = "failed"
	}
	attempt := min(event.Attempt, 10)
	if attempt < 1 {
		attempt = 1
	}
	backoff := r.config.BaseBackoff * time.Duration(1<<(attempt-1))
	return db.New(r.pool).MarkOutboxFailed(ctx, db.MarkOutboxFailedParams{
		ID: event.ID, Status: status,
		Column3:   pgtype.Interval{Microseconds: backoff.Microseconds()},
		LastError: publishErr.Error(), ClaimVersion: event.ClaimVersion,
	})
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
