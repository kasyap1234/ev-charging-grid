package domain

import (
	"encoding/json"
	"time"
	"uuid"
)

type EventType string

const (
	EventPlugHeld         EventType = "PLUG_HELD"
	EventHoldExpired      EventType = "HOLD_EXPIRED"
	EventSessionStarted   EventType = "SESSION_STARTED"
	EventSessionCompleted EventType = "SESSION_COMPLETED"
)

type Event struct {
	ID          string          `json:"id"`
	AggregateID string          `json:"aggregate_id"`
	Type        EventType       `json:"type"`
	Payload     json.RawMessage `json:"payload"`
	CreatedAt   time.Time       `json:"created_at"`
}

type PlugHeldPayload struct {
	SessionID string    `json:"session_id"`
	PlugID    uuid.UUID `json:"plug_id"`
	StationID uuid.UUID `json:"station_id"`
	DriverID  uuid.UUID `json:"driver_id"`
	HeldUntil time.Time `json:"held_until"`
}


