package service

import (
	"context"
	"testing"

	"github.com/feanchy/Order-Service/internal/model"
)

type mockOrderRepository struct {
	order *model.Order
	err   error
}

func (m *mockOrderRepository) GetByID(ctx context.Context, id int) (*model.Order, error) {
	return m.order, m.err
}

func (m *mockOrderRepository) CreateOrder(ctx context.Context) (*model.Order, error) {
	return m.order, m.err
}

func TestOrderService_GetByID(t *testing.T) {
	mockRepo := &mockOrderRepository{
		order: &model.Order{
			ID:     1,
			Status: "created",
		},
	}

	service := NewOrderService(mockRepo)

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
