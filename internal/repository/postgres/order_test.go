package postgres

import (
	"context"
	"testing"

	"github.com/feanchy/Order-Service/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestOrderRepository_GetByID(t *testing.T) {
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, "postgres://postgres:postgres@localhost:5432/orders?sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	query := `
INSERT INTO orders (status)
VALUES ($1)
RETURNING id
`

	order := model.Order{
		Status: "run",
	}

	var orderID int

	err = pool.QueryRow(ctx, query, order.Status).Scan(&orderID)
	if err != nil {
		t.Fatal(err)
	}

	repo := NewOrderRepository(pool)

	got, err := repo.GetByID(ctx, orderID)
	if err != nil {
		t.Fatal(err)
	}

	if got.ID != orderID {
		t.Errorf("expected ID %d, got %d", orderID, got.ID)
	}

	if got.Status != order.Status {
		t.Errorf("expected status %q, got %q", order.Status, got.Status)
	}
}
