package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/gryd-database/platform-poc/configuration"
	"github.com/gryd-database/platform-poc/pkg/transaction"
	"github.com/gryd-database/platform-poc/pkg/transaction/txMock"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"math/big"
	"strings"
	"testing"
)

func TestGetBalance(t *testing.T) {
	var config, _ = configuration.Init()
	var grydContract = config.GRYDContract
	var grydAddress, grydContractABI, _ = setContracts(grydContract.Address, grydContract.ABI)

	t.Parallel()
	ctx := context.Background()
	owner := common.HexToAddress("abcd")
	totalAmount := big.NewInt(100000000000000000)

	t.Run("ok", func(t *testing.T) {
		t.Parallel()

		expectedCallData, err := grydContractABI.Pack("balanceOf", owner)
		if err != nil {
			t.Fatal(err)
		}

		txService := txMock.New(txMock.WithCallFunc(func(ctx context.Context, request *transaction.TxRequest) (result []byte, err error) {
			if *request.To == grydAddress {
				if !bytes.Equal(expectedCallData[:64], request.Data[:64]) {
					return nil, fmt.Errorf("got wrong call data. wanted %x, got %x", expectedCallData, request.Data)
				}
				return totalAmount.FillBytes(make([]byte, 32)), nil
			}
			return nil, errors.New("unexpected call")
		}))

		contract := NewContract(
			&txService, owner, logrus.New(), grydAddress, grydContractABI)

		_, err = contract.GetBalance(ctx)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("unexpected call data", func(t *testing.T) {
		t.Parallel()

		expectedCallData, err := grydContractABI.Pack("balanceOf", owner)
		if err != nil {
			t.Fatal(err)
		}

		txService := txMock.New(txMock.WithCallFunc(func(ctx context.Context, request *transaction.TxRequest) (result []byte, err error) {
			if *request.To == grydAddress {
				if !bytes.Equal(expectedCallData[:64], request.Data[:64]) {
					return nil, fmt.Errorf("got wrong call data. wanted %x, got %x", expectedCallData, request.Data)
				}
			}
			return nil, errors.New("unexpected call")
		}))

		contract := NewContract(
			&txService, owner, logrus.New(), common.HexToAddress("0x000"), abi.ABI{})

		_, err = contract.GetBalance(ctx)
		if err == nil {
			t.Fatal(err)
		}
	})

	t.Run("tx service error", func(t *testing.T) {
		t.Parallel()

		txService := txMock.New(txMock.WithCallFunc(func(ctx context.Context, request *transaction.TxRequest) (result []byte, err error) {
			return nil, errors.New("unexpected call")
		}))

		contract := NewContract(
			&txService, owner, logrus.New(), common.HexToAddress("0x000"), grydContractABI)

		_, err := contract.GetBalance(ctx)
		if err == nil {
			t.Fatal(err)
		}
	})
}

func TestVerifyEvent(t *testing.T) {
	var config, _ = configuration.Init()
	var grydContract = config.GRYDContract
	var grydAddress, grydContractABI, _ = setContracts(grydContract.Address, grydContract.ABI)

	t.Parallel()
	ctx := context.Background()
	owner := common.HexToAddress("abcd")
	txHash := common.HexToHash("xyz0")

	t.Run("ok", func(t *testing.T) {
		t.Parallel()

		txService := txMock.New(
			txMock.WithWaitForReceiptFunc(func(ctx context.Context, trHash common.Hash) (receipt *types.Receipt, err error) {
				if txHash == trHash {
					return &types.Receipt{
						Status: 1,
						Logs: []*types.Log{
							{Topics: []common.Hash{
								grydContractABI.Events["InsertDataSuccess"].ID},
								Address: grydAddress,
							}},
						ContractAddress: grydAddress,
					}, nil
				}
				return nil, errors.New("unknown tx hash")
			}))

		contract := NewContract(
			&txService, owner, logrus.New(), grydAddress, grydContractABI)

		_, err := contract.VerifyEvent(ctx, txHash.String())
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("with incorrect topic", func(t *testing.T) {
		t.Parallel()

		txService := txMock.New(
			txMock.WithWaitForReceiptFunc(func(ctx context.Context, trHash common.Hash) (receipt *types.Receipt, err error) {
				if txHash == trHash {
					return &types.Receipt{
						Status: 1,
						Logs: []*types.Log{
							{Topics: []common.Hash{},
								Address: grydAddress,
							}},
						ContractAddress: grydAddress,
					}, nil
				}
				return nil, errors.New("unknown tx hash")
			}))

		contract := NewContract(
			&txService, owner, logrus.New(), grydAddress, grydContractABI)

		_, err := contract.VerifyEvent(ctx, txHash.String())
		if err == nil {
			t.Fatal(err)
		}
	})

	t.Run("with incorrect receipt", func(t *testing.T) {
		t.Parallel()

		txService := txMock.New(
			txMock.WithWaitForReceiptFunc(func(ctx context.Context, trHash common.Hash) (receipt *types.Receipt, err error) {
				return nil, errors.New("unknown tx hash")
			}))

		contract := NewContract(
			&txService, owner, logrus.New(), grydAddress, abi.ABI{})

		_, err := contract.VerifyEvent(ctx, txHash.String())
		if err == nil {
			t.Fatal(err)
		}
	})
}

func setContracts(address string, jsonABI interface{}) (common.Address, abi.ABI, error) {
	jsonMarshaledABI, err := json.Marshal(jsonABI)
	if err != nil {
		return common.Address{}, abi.ABI{}, fmt.Errorf("unable to marshal json: %w", err)
	}

	jsonToABI, err := abi.JSON(strings.NewReader(string(jsonMarshaledABI)))
	if err != nil {
		return common.Address{}, abi.ABI{}, fmt.Errorf("unable to parse ABI: %w", err)
	}

	return common.HexToAddress(address), jsonToABI, nil
}
