package storage

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
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

	GRYDContractAddress common.Address
	GRYDContractABI     abi.ABI
}

func New(cdb *pgxpool.Pool, logger *logrus.Logger, pool *pgxpool.Pool, grydAddress common.Address, grydABI abi.ABI) *Storage {
	return &Storage{
		cdb:                 cdb,
		logger:              logger,
		pg:                  pool,
		GRYDContractAddress: grydAddress,
		GRYDContractABI:     grydABI,
	}
}

func (s *Storage) Create(ctx context.Context, voStorage *VoStorage) (*DTOStorage, error) {
	resp, err := s.create(ctx, voStorage)
	if err != nil {
		return resp, fmt.Errorf("unable to store tx info in db: %w", err)
	}
	return resp, nil
}
