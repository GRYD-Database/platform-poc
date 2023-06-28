package transaction

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

var (
	_ Backend = (*WrappedBackend)(nil)
)

type WrappedBackend struct {
	backend Backend
}

func NewBackend(backend Backend) *WrappedBackend {
	return &WrappedBackend{
		backend: backend,
	}
}

func (b *WrappedBackend) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	receipt, err := b.backend.TransactionReceipt(ctx, txHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get tx receipt:%w", err)
	}
	return receipt, nil
}

func (b *WrappedBackend) TransactionByHash(ctx context.Context, hash common.Hash) (*types.Transaction, bool, error) {
	tx, isPending, err := b.backend.TransactionByHash(ctx, hash)
	if err != nil {
		return nil, false, fmt.Errorf("failed to get tx by hash:%w", err)
	}
	return tx, isPending, err
}

func (b *WrappedBackend) BlockNumber(ctx context.Context) (uint64, error) {
	blockNumber, err := b.backend.BlockNumber(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get block number:%w", err)
	}
	return blockNumber, nil
}

func (b *WrappedBackend) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	header, err := b.backend.HeaderByNumber(ctx, number)
	if err != nil {
		return nil, fmt.Errorf("failed to get header by number:%w", err)
	}
	return header, nil
}

func (b *WrappedBackend) BalanceAt(ctx context.Context, address common.Address, block *big.Int) (*big.Int, error) {

	balance, err := b.backend.BalanceAt(ctx, address, block)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance at:%w", err)
	}
	return balance, nil
}

func (b *WrappedBackend) NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error) {

	nonce, err := b.backend.NonceAt(ctx, account, blockNumber)
	if err != nil {
		return 0, fmt.Errorf("failed to get nonce at:%w", err)
	}
	return nonce, nil
}

func (b *WrappedBackend) CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error) {

	code, err := b.backend.CodeAt(ctx, contract, blockNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get code at:%w", err)
	}
	return code, nil
}

func (b *WrappedBackend) CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {

	result, err := b.backend.CallContract(ctx, call, blockNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to call contract:%w", err)
	}
	return result, nil
}

func (b *WrappedBackend) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	nonce, err := b.backend.PendingNonceAt(ctx, account)
	if err != nil {
		return 0, fmt.Errorf("failed to get pending nonce:%w", err)
	}
	return nonce, nil
}

func (b *WrappedBackend) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	gasPrice, err := b.backend.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get suggested gas price:%w", err)
	}
	return gasPrice, nil
}

func (b *WrappedBackend) SuggestGasTipCap(ctx context.Context) (*big.Int, error) {
	gasTipCap, err := b.backend.SuggestGasTipCap(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get suggested tip:%w", err)
	}
	return gasTipCap, nil
}

func (b *WrappedBackend) EstimateGas(ctx context.Context, call ethereum.CallMsg) (gas uint64, err error) {
	gas, err = b.backend.EstimateGas(ctx, call)
	if err != nil {
		return 0, fmt.Errorf("failed to get gas estimate:%w", err)
	}
	return gas, nil
}

func (b *WrappedBackend) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	err := b.backend.SendTransaction(ctx, tx)
	if err != nil {
		return fmt.Errorf("failed to send tx:%w", err)
	}
	return nil
}

func (b *WrappedBackend) FilterLogs(ctx context.Context, query ethereum.FilterQuery) ([]types.Log, error) {
	logs, err := b.backend.FilterLogs(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to filter logs:%w", err)
	}
	return logs, nil
}

func (b *WrappedBackend) ChainID(ctx context.Context) (*big.Int, error) {
	chainID, err := b.backend.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get chainID:%w", err)
	}
	return chainID, nil
}

func (b *WrappedBackend) Close() {
	b.backend.Close()
}
