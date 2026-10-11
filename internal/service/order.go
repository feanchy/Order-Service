package service

import (
	"context"
	"log"

	"github.com/feanchy/Order-Service/internal/event"
	"github.com/feanchy/Order-Service/internal/model"
	"github.com/google/uuid"
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

	eventID := uuid.NewString()

	err = s.producer.PublishOrderCreated(ctx, event.OrderCreated{

		EventID: eventID,
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

	processed, err := s.repo.ProcessOrderCreated(ctx, e.EventID, e.OrderID)
	if err != nil {
		return err
	}

	if !processed {
		return nil
	}

	return nil
}
