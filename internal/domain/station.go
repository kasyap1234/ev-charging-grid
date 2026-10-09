package domain

import "uuid"

type ChargingStation struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	TotalPlugs int       `json:"total_plugs"`
}
