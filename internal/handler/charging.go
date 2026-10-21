package handler

import (
	"net/http"
	"time"
	"uuid"

	"github.com/labstack/echo/v4"
)

type holdRequest struct {
	DriverID        string `json:"driver_id"`
	DurationSeconds int    `json:"duration_seconds"`
}

func (h *Handler) HoldPlug(c echo.Context) error {
	plugId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid plug id ")
	}

	var req holdRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	driverID, err := uuid.Parse(req.DriverID)
	session, err := h.repo.HoldPlug(c.Request().Context(), plugId, driverID, time.Duration(req.DurationSeconds)*time.Second)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, session)
}

func (h *Handler) GetAvailability(c echo.Context) error {
	stationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid station id")

	}
	ctx := c.Request().Context()
	count, err := h.repo.GetAvailableCount(ctx, stationID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]int{"available": count})
}
