-- EnqueueOutboxEvent: Persists an outbox event in the caller's transaction.
-- name: EnqueueOutboxEvent :exec
INSERT INTO outbox_events (event_key, topic, message_key, payload)
VALUES ($1, $2, $3, $4::jsonb)
ON CONFLICT (event_key) DO NOTHING;

-- ClaimOutboxEvents: Claims eligible outbox events for publishing.
-- name: ClaimOutboxEvents :many
WITH candidates AS (
    SELECT id
    FROM outbox_events
    WHERE ((status = 'pending' AND next_attempt_at <= current_timestamp)
       OR (status = 'publishing' AND updated_at < current_timestamp - interval '1 minute'))
    ORDER BY id
    FOR UPDATE SKIP LOCKED
    LIMIT $1
)
UPDATE outbox_events e
SET status = 'publishing', attempts = e.attempts + 1,
    claim_version = e.claim_version + 1,
    updated_at = current_timestamp
FROM candidates c
WHERE e.id = c.id
RETURNING e.id, e.topic, e.message_key, e.payload, e.attempts, e.claim_version;

-- ClaimOutboxEventsByTopics: Claims eligible outbox events for selected topics.
-- name: ClaimOutboxEventsByTopics :many
WITH candidates AS (
    SELECT id
    FROM outbox_events
    WHERE ((status = 'pending' AND next_attempt_at <= current_timestamp)
       OR (status = 'publishing' AND updated_at < current_timestamp - interval '1 minute'))
      AND topic = ANY($2::text[])
    ORDER BY id
    FOR UPDATE SKIP LOCKED
    LIMIT $1
)
UPDATE outbox_events e
SET status = 'publishing', attempts = e.attempts + 1,
    claim_version = e.claim_version + 1,
    updated_at = current_timestamp
FROM candidates c
WHERE e.id = c.id
RETURNING e.id, e.topic, e.message_key, e.payload, e.attempts, e.claim_version;

-- MarkOutboxPublished: Marks an event published only for the active claim.
-- name: MarkOutboxPublished :exec
UPDATE outbox_events
SET status = 'published', published_at = current_timestamp,
    updated_at = current_timestamp, last_error = ''
WHERE id = $1 AND status = 'publishing' AND claim_version = $2;

-- MarkOutboxFailed: Records a failed publish only for the active claim.
-- name: MarkOutboxFailed :exec
UPDATE outbox_events
SET status = $2,
    next_attempt_at = current_timestamp + $3::interval,
    last_error = $4,
    updated_at = current_timestamp
WHERE id = $1 AND status = 'publishing' AND claim_version = $5;

-- ReserveConsumerInbox: Claims an event and returns whether this call owns the reservation.
-- name: ReserveConsumerInbox :one
WITH reserved AS (
    INSERT INTO consumer_inbox (
        consumer_name, event_key, topic, partition_id, message_offset,
        status, attempts, reservation_version, lease_until, last_error, processed_at
    )
    VALUES ($1, $2, $3, $4, $5, 'processing', 1, 1,
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
        WHERE ci.consumer_name = $1
          AND ci.event_key = $2
          AND ci.status = 'processed'
    ) AS processed,
    COALESCE(
        (SELECT reservation_version FROM reserved),
        (SELECT ci.reservation_version FROM consumer_inbox ci WHERE ci.consumer_name = $1 AND ci.event_key = $2)
    )::BIGINT AS reservation_version;

-- MarkConsumerInboxProcessed: Completes only the active reservation.
-- name: MarkConsumerInboxProcessed :exec
UPDATE consumer_inbox
SET status = 'processed', processed_at = current_timestamp,
    lease_until = current_timestamp, last_error = ''
WHERE consumer_name = $1 AND event_key = $2
  AND status = 'processing' AND reservation_version = $3;

-- ReleaseConsumerInbox: Releases only the active reservation.
-- name: ReleaseConsumerInbox :exec
UPDATE consumer_inbox
SET status = 'pending', lease_until = current_timestamp,
    last_error = $3
WHERE consumer_name = $1 AND event_key = $2
  AND status = 'processing' AND reservation_version = $4;
