-- name : FetchUnpublishedEvents :many
SELECT id,aggregate_id,event_type,payload,created_at
FROM outbox_events WHERE published_at IS NULL
