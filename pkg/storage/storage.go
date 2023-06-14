package storage

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

type Storage struct {
	cdb    *pgxpool.Pool
	logger *logrus.Logger
}

func New(cdb *pgxpool.Pool, logger *logrus.Logger) *Storage {
	return &Storage{
		cdb:    cdb,
		logger: logger,
	}
}

func (s *Storage) Create() error {
	return nil
}
