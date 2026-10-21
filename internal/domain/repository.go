package domain

import (
	"context"
	"time"
	"uuid"
)

type ChargingRepository interface {
	HoldPlug(ctx context.Context, plugID uuid.UUID, driverID uuid.UUID, duration time.Duration) (*Session, error)
	ExpireHold(ctx context.Context, sessionID uuid.UUID) error
	GetAvailableCount(ctx context.Context, stationID uuid.UUID) (int, error)
}
