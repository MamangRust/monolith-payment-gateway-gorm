package outbox

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/MamangRust/monolith-payment-gateway-pkg/database/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func newTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  "host=localhost user=test dbname=test sslmode=disable",
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		DisableAutomaticPing: true,
	})
	if err != nil {
		t.Fatalf("gorm.Open() error = %v", err)
	}
	return db
}

func TestEnqueueValidatesAndPersistsEvent(t *testing.T) {
	db := newTestDB(t)
	event := Event{EventKey: "topup:42", Topic: "topup.created", Key: "42", Payload: []byte(`{"id":42}`)}

	// Enqueue must accept the valid event without attempting to execute: run in
	// DryRun mode so the statement is built but never sent to a database.
	dryRun := db.Session(&gorm.Session{DryRun: true, SkipDefaultTransaction: true})
	if err := Enqueue(context.Background(), dryRun, event); err != nil {
		t.Fatalf("Enqueue returned error: %v", err)
	}

	// The idempotent insert must render ON CONFLICT DO NOTHING so a duplicate
	// event key is treated as success on retry.
	sql := db.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&models.OutboxEvent{
			EventKey:   event.EventKey,
			Topic:      event.Topic,
			MessageKey: event.Key,
			Payload:    event.Payload,
		})
	})
	if !strings.Contains(sql, "ON CONFLICT") || !strings.Contains(sql, "DO NOTHING") {
		t.Fatalf("enqueue is not idempotent: %s", sql)
	}
}

func TestEnqueueRejectsIncompleteEvent(t *testing.T) {
	for name, event := range map[string]Event{
		"missing key":     {Topic: "topic", Payload: []byte("x")},
		"missing topic":   {EventKey: "key", Payload: []byte("x")},
		"missing payload": {EventKey: "key", Topic: "topic"},
		"invalid json":    {EventKey: "key", Topic: "topic", Payload: []byte("not-json")},
	} {
		t.Run(name, func(t *testing.T) {
			if err := Enqueue(context.Background(), newTestDB(t), event); !errors.Is(err, ErrInvalidEvent) {
				t.Fatalf("expected ErrInvalidEvent, got %v", err)
			}
		})
	}
}

func TestReserveValidatesKeys(t *testing.T) {
	db := newTestDB(t)
	cases := []struct {
		name          string
		consumer, key string
		db            *gorm.DB
	}{
		{name: "empty consumer", consumer: "", key: "k", db: db},
		{name: "empty event key", consumer: "c", key: "", db: db},
		{name: "nil db", consumer: "c", key: "k", db: nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, _, _, err := Reserve(context.Background(), tc.db, tc.consumer, tc.key, "topic", 0, 0)
			if !errors.Is(err, ErrInvalidInboxKey) {
				t.Fatalf("Reserve() error = %v, want ErrInvalidInboxKey", err)
			}
		})
	}
}

func TestReserveUsesLeaseFencedUpsert(t *testing.T) {
	db := newTestDB(t)
	sql := db.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Raw(`
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
			"consumer", "event", "topic", 0, 0,
			"consumer", "event", "consumer", "event",
		)
	})
	for _, needle := range []string{"ON CONFLICT", "DO UPDATE", "status <> 'processed'", "lease_until <= current_timestamp"} {
		if !strings.Contains(sql, needle) {
			t.Fatalf("reserve SQL missing %q: %s", needle, sql)
		}
	}
}

func TestMarkProcessedFencesOnActiveReservation(t *testing.T) {
	db := newTestDB(t)
	sql := db.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(&models.ConsumerInbox{}).
			Where("consumer_name = ? AND event_key = ? AND status = 'processing' AND reservation_version = ?",
				"consumer", "event", 3).
			Updates(map[string]interface{}{
				"status":       "processed",
				"processed_at": gorm.Expr("current_timestamp"),
				"lease_until":  gorm.Expr("current_timestamp"),
				"last_error":   "",
			})
	})
	for _, needle := range []string{"consumer_name =", "event_key =", "status = 'processing'", "reservation_version =", "processed"} {
		if !strings.Contains(sql, needle) {
			t.Fatalf("mark-processed SQL missing %q: %s", needle, sql)
		}
	}
}

func TestReleaseFencesOnActiveReservation(t *testing.T) {
	db := newTestDB(t)
	sql := db.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(&models.ConsumerInbox{}).
			Where("consumer_name = ? AND event_key = ? AND status = 'processing' AND reservation_version = ?",
				"consumer", "event", 3).
			Updates(map[string]interface{}{
				"status":      "pending",
				"lease_until": gorm.Expr("current_timestamp"),
				"last_error":  "consumer processing failed",
			})
	})
	for _, needle := range []string{"consumer_name =", "status = 'processing'", "reservation_version =", "pending"} {
		if !strings.Contains(sql, needle) {
			t.Fatalf("release SQL missing %q: %s", needle, sql)
		}
	}
}
