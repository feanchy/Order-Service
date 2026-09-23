package postgres

import (
	"context"
	"fmt"

	"github.com/feanchy/Order-Service/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgres(ctx context.Context, cfg config.Config) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.DBUser,     // 1
		cfg.DBPassword, // 2
		cfg.DBHost,     // 3
		cfg.DBPort,     // 4
		cfg.DBName,     // 5
	)

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	err = pool.Ping(ctx)
	if err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}
