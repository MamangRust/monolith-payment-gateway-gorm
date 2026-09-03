-- name: InsertCardEventLog :one
INSERT INTO card_event_logs (
    topic,
    event_type,
    card_number,
    reference_id,
    payload
)
VALUES (sqlc.arg(topic), sqlc.arg(event_type), NULLIF(sqlc.arg(card_number)::text, ''), NULLIF(sqlc.arg(reference_id)::text, ''), sqlc.arg(payload)::jsonb)
ON CONFLICT DO NOTHING
RETURNING event_id, topic, event_type, card_number, reference_id, payload, received_at;
