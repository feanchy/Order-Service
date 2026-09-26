package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/feanchy/Order-Service/internal/model"
)

type mockOrderService struct {
	order *model.Order
	err   error
}

func (m *mockOrderService) GetByID(ctx context.Context, id int) (*model.Order, error) {
	return m.order, m.err
}

func (m *mockOrderService) CreateOrder(ctx context.Context) (*model.Order, error) {
	return m.order, m.err
}

func (m *mockOrderService) GetAll(ctx context.Context) ([]*model.Order, error) {
	return []m.order, m.err
}

func TestOrderHandler_GetByID(t *testing.T) {
	mockService := &mockOrderService{
		order: &model.Order{
			ID:     1,
			Status: "created",
		},
	}

	handler := NewOrderHandler(mockService)

	req := httptest.NewRequest(
		http.MethodGet,
		"/orders/1",
		nil,
	)

	req.SetPathValue("id", "1")

	rec := httptest.NewRecorder()

	handler.GetByID(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var order model.Order
	err := json.NewDecoder(rec.Body).Decode(&order)
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

func TestOrderHandler_GetByID_NotFound(t *testing.T) {
	mockService := &mockOrderService{
		err: model.ErrOrderNotFound,
	}

	handler := NewOrderHandler(mockService)

	req := httptest.NewRequest(
		http.MethodGet,
		"/orders/99",
		nil,
	)

	req.SetPathValue("id", "99")

	rec := httptest.NewRecorder()

	handler.GetByID(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestOrderHandler_GetByID_InvalidID(t *testing.T) {
	mockService := &mockOrderService{}

	handler := NewOrderHandler(mockService)

	req := httptest.NewRequest(
		http.MethodGet,
		"/orders/abc",
		nil,
	)
	req.SetPathValue("id", "abc")

	rec := httptest.NewRecorder()

	handler.GetByID(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestOrderHandler_CreateOrder(t *testing.T) {
	mockService := &mockOrderService{
		order: &model.Order{
			ID:     1,
			Status: "created",
		},
	}

	handler := NewOrderHandler(mockService)

	req := httptest.NewRequest(
		http.MethodPost,
		"/orders",
		nil,
	)

	rec := httptest.NewRecorder()

	handler.CreateOrder(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}

	var order model.Order

	err := json.NewDecoder(rec.Body).Decode(&order)
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

func TestOrderHandler_CreateOrder_Error(t *testing.T) {
	mockService := &mockOrderService{
		err: model.ErrInternalServer,
	}

	handler := NewOrderHandler(mockService)

	req := httptest.NewRequest(
		http.MethodPost,
		"/orders",
		nil,
	)

	rec := httptest.NewRecorder()

	handler.CreateOrder(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected code %d, got %d", http.StatusInternalServerError, rec.Code)
	}
}
