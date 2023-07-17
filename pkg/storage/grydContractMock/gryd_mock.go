package grydContractMock

import (
	"context"
	"github.com/gryd-database/platform-poc/cmd/server"
	"github.com/gryd-database/platform-poc/pkg/storage"
	"math/big"
)

type grydContractMock struct {
	getBalance  func(ctx context.Context) (*big.Int, error)
	verifyEvent func(ctx context.Context, hashTx string) (*storage.EventInsertDataSuccess, error)
}

func (g *grydContractMock) VerifyEvent(ctx context.Context, hashTx string) (*storage.EventInsertDataSuccess, error) {
	return g.verifyEvent(ctx, hashTx)
}

func (g *grydContractMock) GetBalance(ctx context.Context) (*big.Int, error) {
	return g.getBalance(ctx)
}

// Option is an option passed to New
type Option func(mock *grydContractMock)

// New creates a new mock
func New(opts ...Option) server.GRYDContract {
	bs := &grydContractMock{}

	for _, o := range opts {
		o(bs)
	}

	return bs
}

func WithGetBalance(f func(ctx context.Context) (*big.Int, error)) Option {
	return func(mock *grydContractMock) {
		mock.getBalance = f
	}
}

func WithVerifyEvent(f func(ctx context.Context, hashTx string) (*storage.EventInsertDataSuccess, error)) Option {
	return func(mock *grydContractMock) {
		mock.verifyEvent = f
	}
}
