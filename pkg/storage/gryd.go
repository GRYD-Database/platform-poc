package storage

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gryd-database/platform-poc/pkg/transaction"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"math/big"
)

var (
	ErrUnprocessableEvent = errors.New("event cannot be processed or does not exist")
)

type Contract struct {
	txService           transaction.Service
	grydContractAddress common.Address
	grydContractABI     abi.ABI
	owner               common.Address
	logger              *logrus.Logger

	Events Events
}

type Events struct {
	TopicInsertDataSuccess common.Hash
}

type EventInsertDataSuccess struct {
	User      common.Address
	QueryType string
}

func NewContract(txService *transaction.Service, owner common.Address, logger *logrus.Logger, grydAddress common.Address, grydABI abi.ABI) *Contract {
	return &Contract{
		txService:           *txService,
		grydContractAddress: grydAddress,
		grydContractABI:     grydABI,
		owner:               owner,
		Events: Events{
			TopicInsertDataSuccess: grydABI.Events["InsertDataSuccess"].ID,
		},
		logger: logger,
	}
}

func (s *Contract) GetBalance(ctx context.Context) (*big.Int, error) {
	callData, err := s.grydContractABI.Pack("balanceOf", s.owner)
	if err != nil {
		return nil, fmt.Errorf("unable to pack callData: %w", err)
	}

	result, err := s.txService.Call(ctx, &transaction.TxRequest{
		To:   &s.grydContractAddress,
		Data: callData,
	})
	if err != nil {
		return nil, fmt.Errorf("err calling tx service: %w", err)
	}

	results, err := s.grydContractABI.Unpack("balanceOf", result)
	if err != nil {
		return nil, fmt.Errorf("unable to unpack callData: %w", err)
	}

	if len(results) == 0 {
		return nil, errors.New("unexpected empty results")
	}

	return abi.ConvertType(results[0], new(big.Int)).(*big.Int), nil
}

func (s *Contract) VerifyEvent(ctx context.Context, hashTx string) (*EventInsertDataSuccess, error) {
	receipt, err := s.txService.WaitForReceipt(ctx, common.HexToHash(hashTx))
	if err != nil {
		return nil, fmt.Errorf("error getting the receipt from tx hash: %s with error: %w", hashTx, err)
	}

	var event EventInsertDataSuccess

	for _, ev := range receipt.Logs {
		if ev.Address == s.grydContractAddress && len(ev.Topics) > 0 && ev.Topics[0] == s.Events.TopicInsertDataSuccess {
			err = transaction.ParseEvent(&s.grydContractABI, "InsertDataSuccess", &event, *ev)
			if err != nil {
				return nil, fmt.Errorf("error parsing event of hash: %s with error: %w", hashTx, err)
			}
		} else {
			return nil, ErrUnprocessableEvent
		}
	}

	return &event, nil
}
