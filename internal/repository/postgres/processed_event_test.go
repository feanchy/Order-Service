package postgres

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestOrderRepository_MarkProcessed(t *testing.T) {
	ctx := context.Background()

	pool, err := pgxpool.New(
		ctx,
		"postgres://postgres:postgres@localhost:5432/orders?sslmode=disable",
	)
	if err != nil {
		t.Fatal(err)
	}

	defer pool.Close()

	repo := NewOrderRepository(pool)

	eventID := "test-event-123"

	_, err = pool.Exec(
		ctx,
		"DELETE FROM processed_events WHERE event_id = $1",
		eventID,
	)
	if err != nil {
		t.Fatal(err)
	}

	got, err := repo.MarkProcessed(ctx, eventID)
	if err != nil {
		t.Fatal(err)
	}

	if !got {
		t.Fatal("expected first call to return true")
	}

	got, err = repo.MarkProcessed(ctx, eventID)
	if err != nil {
		t.Fatal(err)
	}

	if got {
		t.Fatal("expected duplicate event to return false")
	}
}
