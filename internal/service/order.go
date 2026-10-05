package service

import (
	"context"
	"log"

	"github.com/feanchy/Order-Service/internal/event"
	"github.com/feanchy/Order-Service/internal/kafka"
	"github.com/feanchy/Order-Service/internal/model"
)

type OrderService struct {
	repo     OrderRepository
	producer kafka.Producer
}

func NewOrderService(
	repo OrderRepository,
	producer kafka.Producer,
) *OrderService {
	return &OrderService{
		repo:     repo,
		producer: producer,
	}

}

func (s *OrderService) Get(ctx context.Context, id int) (*model.Order, error) {
	return s.repo.Get(ctx, id)
}

func (s *OrderService) CreateOrder(ctx context.Context) (*model.Order, error) {
	log.Println("creating order")

	order, err := s.repo.CreateOrder(ctx)
	if err != nil {
		log.Println("repo error:", err)
		return nil, err
	}

	log.Println("order created:", order.ID)

	err = s.producer.PublishOrderCreated(ctx, event.OrderCreated{
		OrderID: order.ID,
	})
	if err != nil {
		log.Println("kafka error:", err)
		return nil, err
	}

	log.Println("kafka event published")

	return order, nil
}

func (s *OrderService) List(ctx context.Context) ([]*model.Order, error) {
	return s.repo.List(ctx)
}
