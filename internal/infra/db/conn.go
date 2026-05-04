package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func ConnectDB() (*pgxpool.Pool, error) {
	_ = godotenv.Load()

	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		return nil, fmt.Errorf("DATABASE_URL not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	config, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	return pool, nil
}
