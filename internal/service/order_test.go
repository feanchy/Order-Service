package service

import (
	"context"
	"errors"
	"testing"

	"github.com/feanchy/Order-Service/internal/event"
	"github.com/feanchy/Order-Service/internal/model"
)

type mockOrderRepository struct {
	order *model.Order
	err   error
}

type mockProducer struct {
	err   error
	event event.OrderCreated
}

func (m *mockProducer) PublishOrderCreated(ctx context.Context, e event.OrderCreated) error {
	m.event = e
	return m.err
}

func (m *mockOrderRepository) GetByID(ctx context.Context, id int) (*model.Order, error) {
	return m.order, m.err
}

func (m *mockOrderRepository) CreateOrder(ctx context.Context) (*model.Order, error) {
	return m.order, m.err
}

func (m *mockOrderRepository) GetAll(ctx context.Context) ([]*model.Order, error) {
	return []*model.Order{m.order}, m.err
}

func TestOrderService_GetByID(t *testing.T) {
	mockRepo := &mockOrderRepository{
		order: &model.Order{
			ID:     1,
			Status: "created",
		},
	}

	mockProducer := &mockProducer{}

	service := NewOrderService(mockRepo, mockProducer)

	order, err := service.GetByID(context.Background(), 1)

	if err != nil {
		t.Fatal(err)
	}

	if order.ID != 1 {
		t.Errorf("expected ID 1, got %d", order.ID)
	}

	if order.Status != "created" {
		t.Errorf("expected status created, got %s", order.Status)
	}
}

func TestOrderServie_GetByID_NotFound(t *testing.T) {
	mockRepo := &mockOrderRepository{
		err: model.ErrOrderNotFound,
	}

	mockProducer := &mockProducer{}

	service := NewOrderService(mockRepo, mockProducer)

	_, err := service.GetByID(context.Background(), 99)
	if !errors.Is(err, model.ErrOrderNotFound) {
		t.Fatal()
	}
}

func TestOrderService_CreateOrder(t *testing.T) {
	mockRepo := &mockOrderRepository{
		order: &model.Order{
			ID:     1,
			Status: "created",
		},
	}

	mockProducer := &mockProducer{}

	service := NewOrderService(mockRepo, mockProducer)

	order, err := service.CreateOrder(context.Background())
	if err != nil {
		t.Fatal()
	}

	if order.ID != mockRepo.order.ID {
		t.Errorf("expected %d, got %d", mockRepo.order.ID, order.ID)
	}

	if order.Status != mockRepo.order.Status {
		t.Errorf("expected %s, got %s", mockRepo.order.Status, order.Status)
	}

	if mockProducer.event.OrderID != order.ID {
		t.Errorf("expected event OrderID %d, got %d",
			order.ID,
			mockProducer.event.OrderID)
	}
}

func TestOrderService_CreateOrder_Error(t *testing.T) {
	mockRepo := &mockOrderRepository{
		err: errors.New("database error"),
	}

	mockProducer := &mockProducer{}

	service := NewOrderService(mockRepo, mockProducer)

	order, err := service.CreateOrder(context.Background())

	if order != nil {
		t.Fatal("expected order to be nil")
	}

	if err == nil {
		t.Fatal("expected error, got nil")
	}

}
