package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/feanchy/Order-Service/internal/model"
)

type OrderService interface {
	Get(ctx context.Context, id int) (*model.Order, error)
	CreateOrder(ctx context.Context) (*model.Order, error)
	List(ctx context.Context) ([]*model.Order, error)
}

type OrderHandler struct {
	service OrderService
}

func NewOrderHandler(service OrderService) *OrderHandler {
	return &OrderHandler{
		service: service,
	}
}

func (h *OrderHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "invalid order id")
		return
	}

	order, err := h.service.Get(r.Context(), id)

	if errors.Is(err, model.ErrOrderNotFound) {
		writeError(w, http.StatusNotFound, model.ErrOrderNotFound.Error())
		return
	}

	if err != nil {
		log.Println("get order error:", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(order); err != nil {
		http.Error(w, model.ErrInternalServer.Error(), http.StatusInternalServerError)
		return
	}

}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {

	order, err := h.service.CreateOrder(r.Context())
	if err != nil {
		log.Println(err)
		writeError(w, http.StatusInternalServerError, model.ErrInternalServer.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(order); err != nil {
		log.Println(err)
		http.Error(w, model.ErrInternalServer.Error(), http.StatusInternalServerError)
		return
	}

}

func (h *OrderHandler) List(w http.ResponseWriter, r *http.Request) {
	orders, err := h.service.List(r.Context())
	if err != nil {
		log.Println("list orders error:", err)
		writeError(w, http.StatusInternalServerError, model.ErrInternalServer.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(orders); err != nil {
		log.Println("encode orders error:", err)
		return
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}
