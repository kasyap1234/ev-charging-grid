package domain

import "errors"

var (
	ErrStationNotFound = errors.New("station not found")
	ErrSessionNotFound = errors.New("session not found")
	ErrEventNotFound   = errors.New("event not found")
	ErrDriverNotFound  = errors.New("driver not found")
	ErrPlugNotFound    = errors.New("plug not found")

	ErrInvalidDriver       = errors.New("driver id cannot be empty")
	ErrInvalidHoldDuration = errors.New("hold duration must be greater than zero ")
)
