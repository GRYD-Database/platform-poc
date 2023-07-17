package storageMock

import (
	"context"
	"github.com/gryd-database/platform-poc/pkg/storage"
)

type storageMock struct {
	addRecord             func(ctx context.Context, storage *[]storage.InputData) error
	ledger                func(ctx context.Context, wallet, datasetKey string) error
	getWalletByDatasetKey func(ctx context.Context, key string) (*storage.Ledger, error)
	getRecordByID         func(ctx context.Context, id string) (*storage.InputData, error)
}

func (s *storageMock) AddRecord(ctx context.Context, storage *[]storage.InputData) error {
	return s.addRecord(ctx, storage)
}

func (s *storageMock) Ledger(ctx context.Context, wallet, datasetKey string) error {
	return s.ledger(ctx, wallet, datasetKey)
}

func (s *storageMock) GetWalletByDatasetKey(ctx context.Context, key string) (*storage.Ledger, error) {
	return s.getWalletByDatasetKey(ctx, key)
}

func (s *storageMock) GetRecordByID(ctx context.Context, id string) (*storage.InputData, error) {
	return s.getRecordByID(ctx, id)
}

// Option is an option passed to New
type Option func(mock *storageMock)

// New creates a new mock
func New(opts ...Option) storage.OrbitService {
	bs := &storageMock{}

	for _, o := range opts {
		o(bs)
	}

	return bs
}

func WithAddRecord(f func(ctx context.Context, storage *[]storage.InputData) error) Option {
	return func(mock *storageMock) {
		mock.addRecord = f
	}
}

func WithLedger(f func(ctx context.Context, wallet, datasetKey string) error) Option {
	return func(mock *storageMock) {
		mock.ledger = f
	}
}
