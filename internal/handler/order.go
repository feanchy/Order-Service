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
	GetByID(ctx context.Context, id int) (*model.Order, error)
}

type OrderHandler struct {
	service OrderService
}

func NewOrderHandler(service OrderService) *OrderHandler {
	return &OrderHandler{
		service: service,
	}
}

func (h *OrderHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {

		http.Error(w, "invalid order id", http.StatusBadRequest)
		return
	}

	order, err := h.service.GetByID(r.Context(), id)

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
		http.Error(w, "internal server error", http.StatusInternalServerError)
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
