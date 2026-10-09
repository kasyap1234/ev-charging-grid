package domain

import (
	"uuid"
)

type PlugStatus string

const (
	PlugOpen         PlugStatus = "AVAILABLE"
	PlugHeld         PlugStatus = "HELD"
	PlugCharging     PlugStatus = "CHARGING"
	PlugOutOfService PlugStatus = "OUT_OF_SERVICE"
)

type Plug struct {
	ID        uuid.UUID
	StationID uuid.UUID
	Status    PlugStatus
}
