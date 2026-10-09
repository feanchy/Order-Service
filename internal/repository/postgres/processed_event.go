package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func (r *OrderRepository) MarkProcessed(ctx context.Context, eventID string) (bool, error) {
	const query = `
	INSERT INTO processed_events (event_id)
	VALUES ($1)
	ON CONFLICT (event_id) DO NOTHING
	RETURNING event_id
	`

	var id string

	err := r.db.QueryRow(ctx, query, eventID).Scan(&id)
	if err == pgx.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}
