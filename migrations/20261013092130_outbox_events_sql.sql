-- +goose Up
SELECT 'up SQL query';
CREATE TABLE IF NOT EXISTS outbox_events(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    aggregate_id UUID NOT NULL,
    event_type TEXT NOT NULL,
    payload JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
    published_at TIMESTAMPTZ NULL);


CREATE INDEX IF NOT EXISTS idx_outbox_unprocessed
ON outbox_events( created_at ASC)
WHERE published_at IS NULL;
-- +goose Down
SELECT 'down SQL query';
DROP INDEX IF EXISTS idx_outbox_unprocessed;
DROP TABLE IF EXISTS outbox_events;
