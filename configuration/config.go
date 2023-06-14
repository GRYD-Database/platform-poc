package configuration

import (
	"fmt"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	Address  string `env:"ADDRESS"`
	Postgres struct {
		Host       string `env:"DB_HOST"`
		Password   string `env:"DB_PASSWORD"`
		Port       string `env:"DB_PORT"`
		DBName     string `env:"DB_NAME"`
		DBUsername string `env:"DB_USERNAME"`
	}
	CockroachDB struct {
		User     string `env:"CRDB_USER"`
		Password string `env:"CRDB_PASSWORD"`
		Port     string `env:"CRDB_PORT"`
		DBName   string `env:"CRDB_DB"`
		Server   string `env:"CRDB_SERVER"`
		SSLMode  string `env:"CRDB_SSLMODE"`
	}
	Logger struct {
		LogLevel string `env:"LOG_LEVEL"`
		LogEnv   string `env:"LOG_ENV"`
	}
}

func Init() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("error parsing env variables: %w", err)
	}
	return cfg, nil
}
