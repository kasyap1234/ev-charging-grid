package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"time"
	"uuid"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kasyap1234/ev-charging-grid/internal/domain"
	"github.com/kasyap1234/ev-charging-grid/internal/storage/postgres/db"
	"github.com/kasyap1234/ev-charging-grid/pkg"
)

type ChargingRepository struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

var _ domain.ChargingRepository = (*ChargingRepository)(nil)

func NewChargingRepository(pool *pgxpool.Pool) *ChargingRepository {
	return &ChargingRepository{
		pool:    pool,
		queries: db.New(pool),
	}
}

func (r *ChargingRepository) HoldPlug(ctx context.Context, plugID uuid.UUID, driverID uuid.UUID, duration time.Duration) (*domain.Session, error) {
	if duration <= 0 {
		return nil, pkg.ErrInvalidHoldDuration
	}
	if driverID == uuid.Nil() {
		return nil, pkg.ErrInvalidDriver
	}

	plug, err := r.queries.GetPlugByID(ctx, plugID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pkg.ErrPlugNotFound
		}
		return nil, err
	}

	sessionID := uuid.New()
	heldUntil := time.Now().Add(duration)
	payload, err := json.Marshal(domain.PlugHeldPayload{
		SessionID: sessionID.String(),
		PlugID:    plug.ID,
		StationID: plug.StationID,
		DriverID:  driverID,
		HeldUntil: heldUntil,
	})
	if err != nil {
		return nil, err
	}

	row, err := r.queries.HoldPlugAtomic(ctx, db.HoldPlugAtomicParams{
		PlugID:    plugID,
		SessionID: sessionID,
		DriverID:  driverID,
		HeldUntil: heldUntil,
		OutboxID:  uuid.New(),
		Payload:   payload,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pkg.ErrPlugUnavailable
		}
		return nil, err
	}

	return &domain.Session{
		ID:        row.SessionID,
		PlugID:    row.PlugID,
		StationID: row.StationID,
		DriverID:  row.DriverID,
		Status:    row.Status,
		HeldUntil: row.HeldUntil,
		StartedAt: row.StartedAt,
	}, nil
}

func (r *ChargingRepository) ExpireHold(ctx context.Context, sessionID uuid.UUID) error {
	payload, err := json.Marshal(domain.HoldExpiredPayload{
		SessionID: sessionID.String(),
	})
	if err != nil {
		return err
	}

	_, err = r.queries.ExpireHoldAtomic(ctx, db.ExpireHoldAtomicParams{
		SessionID: sessionID,
		OutboxID:  uuid.New(),
		Payload:   payload,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return pkg.ErrSessionNotFound
	}
	return err
}

func (r *ChargingRepository) GetAvailableCount(ctx context.Context, stationID uuid.UUID) (int, error) {
	count, err := r.queries.CountAvailablePlugsByStation(ctx, stationID)
	if err != nil {
		return 0, err
	}
	return int(count), nil
}





