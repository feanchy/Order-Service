package service

import (
	"context"

	"github.com/feanchy/Order-Service/internal/model"
)

type OrderRepository interface {
	GetByID(ctx context.Context, id int) (*model.Order, error)
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
