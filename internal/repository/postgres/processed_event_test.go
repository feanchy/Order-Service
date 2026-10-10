package postgres

import (
	"context"
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

	var orderID int

	err = pool.QueryRow(ctx, `
INSERT INTO orders (status)
VALUES ($1)
RETURNING id
`, model.OrderStatusCreated).Scan(&orderID)
	if err != nil {
		t.Fatal(err)
	}

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

	orderID := 1
	eventID := model.OrderStatusConfirmed

}
