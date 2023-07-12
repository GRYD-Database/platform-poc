package storage

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"time"
)

//nolint:golint,gochecknoglobals,varnamelen
var QB = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

// VoStorage struct contains Wallet and TxHash where Wallet is the address and TxHash is the transaction that was sent over network
type VoStorage struct {
	Wallet     string `json:"wallet"`
	TxHash     string `json:"txHash"`
	DatasetKey string `json:"datasetKey"`
}

type DTOStorage struct {
	ID         uuid.UUID `json:"id"`
	Wallet     string    `json:"wallet"`
	TxHash     string    `json:"txHash"`
	CreatedAt  time.Time `json:"createdAt"`
	DatasetKey string    `json:"datasetKey"`
}

func (s *Storage) Create(ctx context.Context, voStorage *VoStorage) (*DTOStorage, error) {
	resp, err := s.create(ctx, voStorage)
	if err != nil {
		return resp, fmt.Errorf("unable to store tx info in db: %w", err)
	}

	return resp, nil
}

func (s *Storage) create(ctx context.Context, voStorage *VoStorage) (*DTOStorage, error) {
	sqls, args, err := QB.Insert("storage").
		Columns("wallet", "txHash", "datasetKey").
		Values(voStorage.Wallet, voStorage.TxHash, voStorage.DatasetKey).
		Suffix("RETURNING *").
		ToSql()
	if err != nil {
		return &DTOStorage{}, fmt.Errorf("error building query for create storage: %w", err)
	}

	var dtoStorage DTOStorage

	rows, err := s.pg.Query(ctx, sqls, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing query for create storage: %w", err)
	}

	for rows.Next() {
		err := rows.Scan(&dtoStorage.ID, &dtoStorage.Wallet, &dtoStorage.TxHash, &dtoStorage.CreatedAt, &dtoStorage.DatasetKey)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, fmt.Errorf("cannot fetch created rows: %w", err)
			}
			return nil, fmt.Errorf("error scanning for create storage: %w", err)
		}

	}

	return &dtoStorage, nil
}
