package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/feanchy/Order-Service/internal/config"
	"github.com/feanchy/Order-Service/internal/repository/postgres"
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

	// repo := repository.New(pool)
	// service := service.New(repo)
	// handler := handler.New(service)

	server := http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.HTTPPort),
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
		IdleTimeout:  1 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

}
