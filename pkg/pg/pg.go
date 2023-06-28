package pg

import (
	"context"
	"fmt"
	"github.com/gryd-database/platform-poc/configuration"

	"github.com/jackc/pgx/v4/pgxpool"
)

func InitPool(config *configuration.Config) (*pgxpool.Pool, error) {
	args := config.Postgres
	databaseUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", args.DBUsername, args.Password, args.Host, args.Port, args.DBName)

	// this returns connection pool
	pool, err := pgxpool.Connect(context.Background(), databaseUrl)
	if err != nil {
		return nil, fmt.Errorf("error bootstrapping pg: %w", err)
	}

	err = pool.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error pinging pg: %w", err)
	}

	return pool, nil
}
