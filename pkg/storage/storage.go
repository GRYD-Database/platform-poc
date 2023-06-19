package storage

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
	"time"
)

//nolint:golint,gochecknoglobals,varnamelen
var QB = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

// VoStorage struct contains Wallet and TxHash where Wallet is the address and TxHash is the transaction that was sent over network
type VoStorage struct {
	Wallet string `json:"wallet"`
	TxHash string `json:"txHash"`
}

type DTOStorage struct {
	Wallet    string    `json:"wallet"`
	TxHash    string    `json:"txHash"`
	CreatedAt time.Time `json:"createdAt"`
}

type DTODSNInfo struct {
	DSN string `json:"DSN"`
}

type Storage struct {
	cdb    *pgxpool.Pool
	logger *logrus.Logger
	pg     *pgxpool.Pool
}

func New(cdb *pgxpool.Pool, logger *logrus.Logger, pool *pgxpool.Pool) *Storage {
	return &Storage{
		cdb:    cdb,
		logger: logger,
		pg:     pool,
	}
}

func (s *Storage) Create(ctx context.Context, voStorage *VoStorage) (*DTOStorage, error) {
	resp, err := s.create(ctx, voStorage)
	if err != nil {
		return resp, fmt.Errorf("unable to store tx info in db: %w", err)
	}
	return resp, nil
}
