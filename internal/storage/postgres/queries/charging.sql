-- ============================================================================
-- 1. ATOMIC HOLD RESERVATION (Single-Hop CTE)
-- Invariant: Locks plug, marks it HELD, creates session, and writes Outbox event.
-- If plug is already HELD or does not exist, 0 rows are returned (pgx.ErrNoRows).
-- ============================================================================

-- name: HoldPlugAtomic :one
WITH locked_plug AS (
    SELECT id, station_id
    FROM charging_plugs
    WHERE id = @plug_id AND status = 'AVAILABLE'
    FOR UPDATE
),
updated_plug AS (
    UPDATE charging_plugs
    SET status = 'HELD'
    WHERE id IN (SELECT id FROM locked_plug)
    RETURNING id, station_id
),
new_session AS (
    INSERT INTO charging_sessions (
        id,
        plug_id,
        driver_id,
        status,
        held_until,
        started_at
    )
    SELECT
        @session_id,
        id,
        @driver_id,
        'HELD',
        @held_until,
        NOW()
    FROM updated_plug
    RETURNING id, plug_id, driver_id, status, held_until, started_at
),
inserted_outbox AS (
    INSERT INTO outbox_events (
        id,
        aggregate_id,
        event_type,
        payload,
        created_at
    )
    SELECT
        @outbox_id,
        station_id,
        'PLUG_HELD',
        @payload,
        NOW()
    FROM updated_plug
)
SELECT 
    s.id AS session_id,
    s.plug_id,
    s.driver_id,
    s.status,
    s.held_until,
    s.started_at,
    p.station_id
FROM new_session s
JOIN updated_plug p ON s.plug_id = p.id;

-- ============================================================================
-- 2. ATOMIC EXPIRATION (Sweeper Cleanup CTE)
-- Invariant: Releases an expired hold back to AVAILABLE and publishes HOLD_EXPIRED.
-- ============================================================================

-- name: ExpireHoldAtomic :one
WITH target_session AS (
    SELECT id, plug_id
    FROM charging_sessions
    WHERE id = @session_id AND status = 'HELD'
    FOR UPDATE
),
updated_session AS (
    UPDATE charging_sessions
    SET status = 'EXPIRED', ended_at = NOW()
    WHERE id IN (SELECT id FROM target_session)
    RETURNING id, plug_id
),
freed_plug AS (
    UPDATE charging_plugs
    SET status = 'AVAILABLE'
    WHERE id IN (SELECT plug_id FROM updated_session)
    RETURNING id, station_id
),
inserted_outbox AS (
    INSERT INTO outbox_events (
        id,
        aggregate_id,
        event_type,
        payload,
        created_at
    )
    SELECT
        @outbox_id,
        station_id,
        'HOLD_EXPIRED',
        @payload,
        NOW()
    FROM freed_plug
)
SELECT 
    s.id AS session_id,
    p.id AS plug_id,
    p.station_id
FROM updated_session s
JOIN freed_plug p ON s.plug_id = p.id;

-- ============================================================================
-- 3. READ QUERIES (Cache Hydration & Verification)
-- ============================================================================

-- name: GetAvailablePlugsByStation :many
SELECT id, station_id, status
FROM charging_plugs
WHERE station_id = @station_id AND status = 'AVAILABLE';

-- name: CountAvailablePlugsByStation :one
SELECT COUNT(*)::bigint AS count
FROM charging_plugs
WHERE station_id = @station_id AND status = 'AVAILABLE';