-- name: CreateResetToken :one
WITH inserted_token AS (
    INSERT INTO reset_tokens (user_id, token, expiry_date)
    VALUES ($1, $2, $3)
    RETURNING user_id, token, expiry_date
), queued_event AS (
    INSERT INTO outbox_events (event_key, topic, message_key, payload)
    SELECT
        'auth-forgot-password:' || it.user_id::text || ':' || it.token,
        'email-service-topic-auth-forgot-password',
        it.user_id::text,
        jsonb_build_object(
            'event_id', gen_random_uuid()::text,
            'schema_version', 1,
            'event_type', 'auth.forgot_password',
            'occurred_at', now(),
            'email', u.email,
            'subject', 'Password Reset Request',
            'body', 'Click the reset link to change your password: https://sanedge.example.com/reset-password?token=' || it.token
        )
    FROM inserted_token it
    JOIN users u ON u.user_id = it.user_id AND u.deleted_at IS NULL
    ON CONFLICT (event_key) DO NOTHING
    RETURNING event_key
)
SELECT user_id, token, expiry_date
FROM inserted_token;

-- name: DeleteResetToken :exec
DELETE FROM reset_tokens WHERE user_id = $1;

-- name: GetResetToken :one
SELECT user_id, token, expiry_date
FROM reset_tokens
WHERE
    token = $1;