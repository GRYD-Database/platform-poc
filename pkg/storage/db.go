package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

func (s *Storage) create(ctx context.Context, voStorage *VoStorage) (*DTOStorage, error) {
	sqls, args, err := QB.Insert("storage").
		Columns("wallet", "txHash").
		Values(voStorage.Wallet, voStorage.TxHash).
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
		err := rows.Scan(&dtoStorage.Wallet, &dtoStorage.TxHash, &dtoStorage.CreatedAt)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, nil
			}
			return nil, fmt.Errorf("error scanning for create storage: %w", err)
		}

	}

	return &dtoStorage, nil
}
