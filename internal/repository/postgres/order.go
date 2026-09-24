package postgres

import (
	"context"

	"github.com/feanchy/Order-Service/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

func (r *OrderRepository) GetByID(ctx context.Context, id int) (*model.Order, error) {
	const query = `
	SELECT id, status 
	FROM orders
	WHERE id = $1`

	var order model.Order

	err := r.db.QueryRow(ctx, query, id).Scan(
		&order.ID,
		&order.Status,
	)
	if err != nil {
		return nil, err
	}

	return &order, nil

}

func (r *OrderRepository) Create(ctx context.Context, order model.Order) (int, error) {
	const query = `
        INSERT INTO orders (status)
        VALUES ($1)
        RETURNING id
    `

	var id int

	err := r.db.QueryRow(ctx, query, order.Status).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *OrderRepository) UpdateStatus(ctx context.Context, order model.Order) error {
	const query = `
	UPDATE orders
	SET status = $1
	WHERE id = $2
	`

	_, err := r.db.Exec(ctx, query, order.Status, order.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *OrderRepository) Delete(ctx context.Context, id int) error {
	const query = `
	DELETE FROM orders
	WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	return nil

}
