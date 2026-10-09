package domain

import (
	"time"
	"uuid"
)

type Session struct {
	ID        uuid.UUID
	PlugID    uuid.UUID
	StationID uuid.UUID
	DriverID  uuid.UUID
	Status    string
	HeldUntil time.Time
	StartedAt time.Time
}
