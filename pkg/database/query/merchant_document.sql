-- name: GetMerchantDocuments :many
SELECT
    document_id,
    merchant_id,
    document_type,
    document_url,
    status,
    note,
    uploaded_at,
    created_at,
    updated_at,
    COUNT(*) OVER () AS total_count
FROM merchant_documents
WHERE
    deleted_at IS NULL
    AND (
        $1::TEXT IS NULL
        OR document_type ILIKE '%' || $1 || '%'
        OR status ILIKE '%' || $1 || '%'
        OR note ILIKE '%' || $1 || '%'
    )
ORDER BY document_id
LIMIT $2
OFFSET
    $3;

-- name: GetActiveMerchantDocuments :many
SELECT
    document_id,
    merchant_id,
    document_type,
    document_url,
    status,
    note,
    uploaded_at,
    created_at,
    updated_at,
    deleted_at,
    COUNT(*) OVER () AS total_count
FROM merchant_documents
WHERE
    deleted_at IS NULL
    AND status != 'deleted'
    AND (
        $1::TEXT IS NULL
        OR document_type ILIKE '%' || $1 || '%'
        OR status ILIKE '%' || $1 || '%'
        OR note ILIKE '%' || $1 || '%'
    )
ORDER BY document_id
LIMIT $2
OFFSET
    $3;

-- name: GetTrashedMerchantDocuments :many
SELECT
    document_id,
    merchant_id,
    document_type,
    document_url,
    status,
    note,
    uploaded_at,
    created_at,
    updated_at,
    deleted_at,
    COUNT(*) OVER () AS total_count
FROM merchant_documents
WHERE
    deleted_at IS NOT NULL
    AND (
        $1::TEXT IS NULL
        OR document_type ILIKE '%' || $1 || '%'
        OR status ILIKE '%' || $1 || '%'
        OR note ILIKE '%' || $1 || '%'
    )
ORDER BY document_id
LIMIT $2
OFFSET
    $3;

-- name: GetMerchantDocument :one
SELECT
    document_id,
    merchant_id,
    document_type,
    document_url,
    status,
    note,
    uploaded_at,
    created_at,
    updated_at
FROM merchant_documents
WHERE
    document_id = $1
    AND deleted_at IS NULL;

-- name: CreateMerchantDocument :one
WITH valid_merchant AS (
    SELECT m.merchant_id, u.email
    FROM merchants m
    JOIN users u ON u.user_id = m.user_id AND u.deleted_at IS NULL
    WHERE m.merchant_id = $1 AND m.deleted_at IS NULL
), inserted_document AS (
    INSERT INTO merchant_documents (
        merchant_id, document_type, document_url, status, note, uploaded_at, updated_at
    )
    SELECT merchant_id, $2, $3, $4, $5, current_timestamp, current_timestamp
    FROM valid_merchant
    RETURNING document_id, merchant_id, document_type, document_url, status, note, uploaded_at, created_at, updated_at
), queued_event AS (
    INSERT INTO outbox_events (event_key, topic, message_key, payload)
    SELECT
        'merchant-document-create:' || document_id::text,
        'email-service-topic-merchant-document-create',
        document_id::text,
        jsonb_build_object(
            'event_id', gen_random_uuid()::text,
            'schema_version', 1,
            'event_type', 'merchant_document.created',
            'occurred_at', now(),
            'email', u.email,
            'subject', 'Merchant Verification Pending - Action Required',
            'body', 'Your merchant document has been submitted and is pending review.'
        )
    FROM inserted_document d
    JOIN valid_merchant u ON u.merchant_id = d.merchant_id
    ON CONFLICT (event_key) DO NOTHING
    RETURNING event_key
)
SELECT document_id, merchant_id, document_type, document_url, status, note, uploaded_at, created_at, updated_at
FROM inserted_document;

-- name: UpdateMerchantDocument :one
UPDATE merchant_documents
SET
    document_type = $2,
    document_url = $3,
    status = $4,
    note = $5,
    updated_at = current_timestamp
WHERE
    document_id = $1
    AND deleted_at IS NULL
RETURNING
    document_id,
    merchant_id,
    document_type,
    document_url,
    status,
    note,
    uploaded_at,
    created_at,
    updated_at;

-- name: UpdateMerchantDocumentStatus :one
WITH valid_merchant AS (
    SELECT m.merchant_id
    FROM merchants m
    JOIN users u ON u.user_id = m.user_id AND u.deleted_at IS NULL
    WHERE m.merchant_id = (
        SELECT md.merchant_id FROM merchant_documents md
        WHERE md.document_id = $1 AND md.deleted_at IS NULL
    ) AND m.deleted_at IS NULL
), updated_document AS (
    UPDATE merchant_documents
    SET status = $2, note = $3, event_version = event_version + 1, updated_at = current_timestamp
    WHERE merchant_documents.document_id = $1
      AND merchant_documents.deleted_at IS NULL
      AND merchant_documents.merchant_id IN (SELECT valid_merchant.merchant_id FROM valid_merchant)
    RETURNING document_id, merchant_id, document_type, document_url, status, note, uploaded_at, created_at, updated_at, event_version
), queued_event AS (
    INSERT INTO outbox_events (event_key, topic, message_key, payload)
    SELECT
        'merchant-document-status:' || d.document_id::text || ':' || d.event_version::text,
        'email-service-topic-merchant-document-update-status',
        d.document_id::text,
        jsonb_build_object(
            'event_id', gen_random_uuid()::text,
            'schema_version', 1,
            'event_type', 'merchant_document.status_updated',
            'occurred_at', now(),
            'email', u.email,
            'subject', 'Merchant Document Status: ' || initcap(d.status),
            'body', 'Your merchant document status has been updated to ' || d.status || '.' || CASE WHEN d.note IS NULL OR d.note = '' THEN '' ELSE ' Reviewer note: ' || d.note END,
            'event_version', d.event_version
        )
    FROM updated_document d
    JOIN merchants m ON m.merchant_id = d.merchant_id AND m.deleted_at IS NULL
    JOIN users u ON u.user_id = m.user_id AND u.deleted_at IS NULL
    ON CONFLICT (event_key) DO NOTHING
    RETURNING event_key
)
SELECT d.document_id, d.merchant_id, d.document_type, d.document_url, d.status, d.note, d.uploaded_at, d.created_at, d.updated_at
FROM updated_document d;

-- name: TrashMerchantDocument :one
UPDATE merchant_documents
SET
    deleted_at = current_timestamp
WHERE
    document_id = $1
    AND deleted_at IS NULL
RETURNING
    document_id,
    merchant_id,
    document_type,
    document_url,
    status,
    note,
    uploaded_at,
    created_at,
    updated_at,
    deleted_at;

-- name: RestoreMerchantDocument :one
UPDATE merchant_documents
SET
    deleted_at = NULL
WHERE
    document_id = $1
    AND deleted_at IS NOT NULL
RETURNING
    document_id,
    merchant_id,
    document_type,
    document_url,
    status,
    note,
    uploaded_at,
    created_at,
    updated_at,
    deleted_at;

-- name: DeleteMerchantDocumentPermanently :exec
DELETE FROM merchant_documents
WHERE
    document_id = $1
    AND deleted_at IS NOT NULL;

-- name: RestoreAllMerchantDocuments :exec
UPDATE merchant_documents
SET
    deleted_at = NULL
WHERE
    deleted_at IS NOT NULL;

-- name: DeleteAllPermanentMerchantDocuments :exec
DELETE FROM merchant_documents WHERE deleted_at IS NOT NULL;