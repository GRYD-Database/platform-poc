package storage

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gofrs/uuid"
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
	ID        uuid.UUID `json:"id"`
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

	owner common.Address
}

func New(owner common.Address, logger *logrus.Logger, pool *pgxpool.Pool) *Storage {
	return &Storage{
		logger: logger,
		pg:     pool,
		owner:  owner,
	}
}

func (s *Storage) AssignStorage() (string, error) {
	return "", nil
}
