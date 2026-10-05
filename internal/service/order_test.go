package service

import (
	"context"
	"errors"
	"testing"

	"github.com/feanchy/Order-Service/internal/event"
	"github.com/feanchy/Order-Service/internal/model"
	"go.uber.org/mock/gomock"
)

func TestOrderService_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockOrderRepository(ctrl)
	producer := NewMockOrderProducer(ctrl)

	order := &model.Order{ID: 1, Status: "created"}

	repo.EXPECT().Get(gomock.Any(), 1).Return(order, nil)

	service := NewOrderService(repo, producer)

	result, err := service.Get(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("expected order, got nil")
	}

	if result.ID != order.ID {
		t.Fatalf("expected ID %d, got %d", order.ID, result.ID)
	}
}

func TestOrderServie_Get_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockOrderRepository(ctrl)
	producer := NewMockOrderProducer(ctrl)

	repo.EXPECT().Get(gomock.Any(), 99).Return(nil, model.ErrOrderNotFound)

	service := NewOrderService(repo, producer)

	_, err := service.Get(context.Background(), 99)
	if !errors.Is(err, model.ErrOrderNotFound) {
		t.Fatalf("expected ErrOrderNotFound, got %v", err)
	}
}

func TestOrderService_CreateOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockOrderRepository(ctrl)
	producer := NewMockOrderProducer(ctrl)

	order := &model.Order{ID: 1, Status: "created"}

	repo.EXPECT().CreateOrder(gomock.Any()).Return(order, nil)
	producer.EXPECT().PublishOrderCreated(gomock.Any(), event.OrderCreated{OrderID: order.ID}).Return(nil)

	service := NewOrderService(repo, producer)

	created, err := service.CreateOrder(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if created.ID != order.ID {
		t.Fatalf("expected ID %d, got %d", order.ID, created.ID)
	}

	if created.Status != order.Status {
		t.Fatalf("expected status %s, got %s", order.Status, created.Status)
	}
}

func TestOrderService_CreateOrder_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockOrderRepository(ctrl)
	producer := NewMockOrderProducer(ctrl)

	repo.EXPECT().CreateOrder(gomock.Any()).Return(nil, errors.New("database error"))

	service := NewOrderService(repo, producer)

	created, err := service.CreateOrder(context.Background())
	if created != nil {
		t.Fatal("expected created order to be nil")
	}

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
