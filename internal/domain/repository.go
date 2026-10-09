package domain

import (
	"context"
	"time"
	"uuid"
)

type ChargingRepository interface {
	HoldPlug(ctx context.Context, plugID uuid.UUID, driverID string, duration time.Duration) (*Session, error)
	GetAvailableCount(ctx context.Context, stationID uuid.UUID) (int, error)
}
