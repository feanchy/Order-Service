package service

import (
	"context"

	"github.com/feanchy/Order-Service/internal/model"
)

type OrderRepository interface {
	GetByID(ctx context.Context, id int) (*model.Order, error)
	CreateOrder(ctx context.Context, status string) (*model.Order, error)
}

type OrderService struct {
	repo OrderRepository
}

func NewOrderService(repo OrderRepository) *OrderService {
	return &OrderService{
		repo: repo,
	}

}

func (s *OrderService) GetByID(ctx context.Context, id int) (*model.Order, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *OrderService) CreateOrder(ctx context.Context) (*model.Order, error) {

	return nil, model.ErrInvalidStatus
}
