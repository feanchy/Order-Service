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
	query := `SELECT id, status FROM orders WHERE id = $1`

	var order model.Order

	err := r.db.QueryRow(ctx, query, id).Scan(&order.ID, &order.Status)
	if err != nil {
		return nil, err
	}

	return &order, nil

}

func (r *OrderRepository) Create(ctx context.Context, id int) (int, error) {

}

func (r *OrderRepository) UpdateStatus(ctx context.Context, id int) (int, error) {

}

func (r *OrderRepository) Delete(ctx context.Context, id int) (int, error) {

}
