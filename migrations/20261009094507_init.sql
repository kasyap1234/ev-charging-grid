-- +goose Up
SELECT 'up SQL query';
CREATE TABLE IF NOT EXISTS charging_stations(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), 
    name TEXT NOT NULL, 
    total_plugs INT NOT NULL, 
);

CREATE TABLE IF NOT EXISTS charging_plugs(
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
station_id UUID NOT NULL REFERENCES charging_stations(id) ON DELETE CASCADE , 
status TEXT NOT NULL DEFAULT 'AVAILABLE',
CONSTRAINT chk_plug_status CHECK (status IN ('AVAILABLE','HELD','CHARGING','OUT_OF_SERVICE'))
);

CREATE TABLE IF NOT EXISTS charging_sessions(
id UUID PRIMARY KEY DEFAULT gen_random_uuid(), 
plug_id UUID NOT NULL REFERENCES charging_plugs(id) ON DELETE CASCADE, 
driver_id UUID NOT NULL, 
status text NOT NULL DEFAULT 'HELD', 
held_until TIMESTAMPTZ NOT NULL, 
started_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), 
ended_at TIMESTAMPTZ NULL, 
kwh_consumed NUMERIC(8,2) DEFAULT 0.00, 
CONSTRAINT chk_session_status CHECK (status IN ('HELD', 'CHARGING','COMPLETED','EXPIRED'))
);
-- +goose Down
SELECT 'down SQL query';
DROP TABLE IF NOT EXISTS charging_sessions; 
DROP TABLE IF NOT EXISTS charging_plugs; 
DROP TABLE IF NOT EXISTS charging_stations; 

