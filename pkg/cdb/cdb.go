package cdb

import (
	"context"
	"fmt"
	"github.com/gryd-database/platform-poc/configuration"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

func Init(config *configuration.Config) (*pgxpool.Pool, error) {
	args := config.CockroachDB

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", args.User, args.Password, args.Server, args.Port, args.DBName, args.SSLMode)

	ctx := context.Background()
	conn, err := pgxpool.Connect(context.Background(), dsn)

	var now time.Time
	err = conn.QueryRow(ctx, "SELECT NOW()").Scan(&now)
	if err != nil {
		return nil, fmt.Errorf("failed to execute init query: %w", err)
	}

	return conn, nil
}
