package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/feanchy/Order-Service/internal/config"
	"github.com/feanchy/Order-Service/internal/handler"
	"github.com/feanchy/Order-Service/internal/repository/postgres"
	"github.com/feanchy/Order-Service/internal/service"
)

func main() {
	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := postgres.NewPostgres(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}

	defer pool.Close()

	repo := postgres.NewOrderRepository(pool)
	OrderService := service.NewOrderService(repo)
	OrderHandler := handler.NewOrderHandler(OrderService)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /orders/{id}", OrderHandler.GetByID)

	server := http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.HTTPPort),
		Handler:      mux,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
		IdleTimeout:  1 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

}
