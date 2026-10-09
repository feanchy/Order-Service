package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

func (r *OrderRepository) ProcessOrderCreated(ctx context.Context, eventID string) (bool, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return false, err
	}
	defer tx.Rollback(ctx)

	const query = `
	INSERT INTO processed_events (event_id)
	VALUES ($1)
	ON CONFLICT (event_id) DO NOTHING
	RETURNING event_id
	`

	var id string

	err = tx.QueryRow(ctx, query, eventID).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if err := tx.Commit(ctx); err != nil {
		return false, err
	}

	return true, nil
}
