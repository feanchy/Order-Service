package postgres

import (
	"context"
	"errors"
	"testing"

	"github.com/feanchy/Order-Service/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestOrderRepository_ProcessOrderCreated(t *testing.T) {
	ctx := context.Background()

	pool, err := pgxpool.New(
		ctx,
		"postgres://postgres:postgres@localhost:5432/orders?sslmode=disable",
	)
	if err != nil {
		t.Fatal(err)
	}

	defer pool.Close()

	eventID := "test-event-123"

	repo := NewOrderRepository(pool)

	var orderID int

	err = pool.QueryRow(ctx, `
INSERT INTO orders (status)
VALUES ($1)
RETURNING id
`, model.OrderStatusCreated).Scan(&orderID)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		_, err := pool.Exec(ctx,
			"DELETE FROM processed_events WHERE event_id = $1",
			eventID,
		)
		if err != nil {
			t.Errorf("cleanup processed event: %v", err)
		}

		_, err = pool.Exec(ctx,
			"DELETE FROM orders WHERE id = $1",
			orderID,
		)
		if err != nil {
			t.Errorf("cleanup order: %v", err)
		}
	}()

	processed, err := repo.ProcessOrderCreated(ctx, eventID, orderID)
	if err != nil {
		t.Fatal(err)
	}
	if !processed {
		t.Fatal("expected first call to return true")
	}

	processed, err = repo.ProcessOrderCreated(ctx, eventID, orderID)
	if err != nil {
		t.Fatal(err)
	}

	if processed {
		t.Fatal("expected duplicate event to return false")
	}

	var status string
	err = pool.QueryRow(
		ctx,
		"SELECT status FROM orders WHERE id = $1",
		orderID,
	).Scan(&status)

	if err != nil {
		t.Fatal(err)
	}
	if status != model.OrderStatusConfirmed {
		t.Fatalf("expected status %q, got %q",
			model.OrderStatusConfirmed, status)
	}
}

func TestOrderRepository_NotFoundOrder(t *testing.T) {
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

	orderID := -1
	eventID := "test-event-not-found"

	_, err = repo.ProcessOrderCreated(ctx, eventID, orderID)
	if !errors.Is(err, model.ErrOrderNotFound) {
		t.Fatalf("expected ErrOrderNotFound, got %v", err)
	}

	defer func() {
		_, err := pool.Exec(
			ctx,
			"DELETE FROM processed_events WHERE event_id = $1",
			eventID,
		)

		if err != nil {
			t.Errorf("cleanup processed event: %v", err)
		}
	}()

	const query = `
SELECT COUNT(*) FROM processed_events WHERE event_id = $1`

	var countRows int

	err = pool.QueryRow(ctx, query, eventID).Scan(&countRows)
	if err != nil {
		t.Fatal(err)
	}

	if countRows != 0 {
		t.Fatalf("expected 0, got %d", countRows)
	}

}
