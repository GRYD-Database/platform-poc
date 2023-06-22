package storage

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/gryd-database/platform-poc/pkg/transaction"
	"github.com/pkg/errors"
	"math/big"
)

func (s *Storage) GetBalance(ctx context.Context) (*big.Int, error) {
	return s.getBalance(ctx)
}

func (s *Storage) getBalance(ctx context.Context) (*big.Int, error) {
	callData, err := s.grydContractABI.Pack("balanceOf", s.owner)
	if err != nil {
		return nil, err
	}

	result, err := s.txService.Call(ctx, &transaction.TxRequest{
		To:   &s.grydContractAddress,
		Data: callData,
	})
	if err != nil {
		return nil, err
	}

	results, err := s.grydContractABI.Unpack("balanceOf", result)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, errors.New("unexpected empty results")
	}

	return abi.ConvertType(results[0], new(big.Int)).(*big.Int), nil
}

func (s *Storage) VerifyEvent(ctx context.Context) error {
	return nil
}
