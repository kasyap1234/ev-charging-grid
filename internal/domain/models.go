package domain

import (
	"context"
	"time"
	"uuid"
)

type PlugStatus string

const (
	PlugOpen         PlugStatus = "AVAILABLE"
	PlugHeld         PlugStatus = "HELD"
	PlugCharging     PlugStatus = "CHARGING"
	PlugOutOfService PlugStatus = "OUT_OF_SERVICE"
)

type ChargingRepository interface {
	HoldPlug(ctx context.Context, plugID uuid.UUID, driverID string, duration time.Duration) (*Session, error)
	GetAvailableCount(ctx context.Context, stationID uuid.UUID) (int, error)
}

type Plug struct {
	ID        uuid.UUID
	StationID uuid.UUID
	Status    PlugStatus
}
