package service

import (
	"context"
	"log"

	"github.com/feanchy/Order-Service/internal/event"
	"github.com/feanchy/Order-Service/internal/model"
)

type OrderService struct {
	repo     OrderRepository
	producer OrderProducer
}

func NewOrderService(
	repo OrderRepository,
	producer OrderProducer,
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

	log.Println("Create Order: order created:", order.ID)

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

func (s *OrderService) HandleOrderCreated(
	ctx context.Context,
	e event.OrderCreated,
) error {

	log.Println("Consumer: order created:", e.OrderID)

	return nil
}
