package handler

import (
	"github.com/kasyap1234/ev-charging-grid/internal/domain"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	repo domain.ChargingRepository
}


func NewHandler(repo domain.ChargingRepository)*Handler {
	return &Handler{
		repo: repo , 
	}
}


func(h*Handler)Register(e*echo.Echo){
	e.POST("/plugs/:id/hold",h.HoldPlug)
	e.GET("/stations/:id/availability",h.GetAvailability)
}