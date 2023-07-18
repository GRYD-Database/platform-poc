package dbMock

import (
	"context"
	"github.com/gryd-database/platform-poc/pkg/storage"
)

type dbMock struct {
	create func(ctx context.Context, voStorage *storage.VoStorage) (*storage.DTOStorage, error)
}

func (s *dbMock) Create(ctx context.Context, voStorage *storage.VoStorage) (*storage.DTOStorage, error) {
	return s.create(ctx, voStorage)
}

type Option func(mock *dbMock)

// New creates a new mock
func New(opts ...Option) storage.DBService {
	bs := &dbMock{}

	for _, o := range opts {
		o(bs)
	}

	return bs
}

func WithCreate(f func(ctx context.Context, voStorage *storage.VoStorage) (*storage.DTOStorage, error)) Option {
	return func(mock *dbMock) {
		mock.create = f
	}
}
