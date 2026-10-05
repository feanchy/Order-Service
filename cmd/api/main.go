package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/feanchy/Order-Service/internal/config"
	"github.com/feanchy/Order-Service/internal/handler"
	"github.com/feanchy/Order-Service/internal/kafka"
	"github.com/feanchy/Order-Service/internal/repository/postgres"
	"github.com/feanchy/Order-Service/internal/service"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("failed to load .env")
	}

	cfg := config.Load()

	fmt.Println("KAFKA_BROKER:", cfg.KafkaBroker)
	fmt.Println("KAFKA_TOPIC:", cfg.KafkaTopic)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := postgres.NewPostgres(ctx, *cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	repo := postgres.NewOrderRepository(pool)

	producer := kafka.NewKafkaProducer(
		[]string{cfg.KafkaBroker},
		cfg.KafkaTopic,
	)

	defer producer.Close()

	orderService := service.NewOrderService(repo, producer)
	orderHandler := handler.NewOrderHandler(orderService)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /orders/{id}", orderHandler.GetByID)
	mux.HandleFunc("POST /orders", orderHandler.CreateOrder)
	mux.HandleFunc("GET /orders", orderHandler.GetAll)

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
