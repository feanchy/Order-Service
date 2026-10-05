package service

import (
	"context"

	"github.com/feanchy/Order-Service/internal/event"
	"github.com/feanchy/Order-Service/internal/model"
)

type OrderRepository interface {
	Get(ctx context.Context, id int) (*model.Order, error)
	CreateOrder(ctx context.Context) (*model.Order, error)
	List(ctx context.Context) ([]*model.Order, error)
}

type OrderProducer interface {
	PublishOrderCreated(ctx context.Context, event event.OrderCreated) error
}

//go:generate mockgen -source=deps.go -package=service -destination=deps_mock.go
