package storage

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gryd-database/platform-poc/pkg/transaction"
	"github.com/pkg/errors"
	"math/big"
)

type eventBuyStorage struct {
	Buyer    common.Address
	UserName string
	Size     *big.Int
}

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

//func (s *Storage) VerifyEvent(ctx context.Context, hashTx string) (bool, error) {
//	hash := common.HexToHash(hashTx)
//	query := ethereum.FilterQuery{
//		Addresses: []common.Address{s.grydContractAddress},
//		BlockHash: &hash,
//		Topics: [][]common.Hash{
//			{
//				s.Events.storageBoughtTopic,
//			},
//		},
//	}
//	logs, err := s.txService.FilterLogs(ctx, query)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	for {
//
//	}
//	return false, nil
//}

func (s *Storage) VerifyEvent(ctx context.Context, hashTx string) (bool, error) {
	receipt, err := s.txService.WaitForReceipt(ctx, common.HexToHash(hashTx))
	if err != nil {
		return false, fmt.Errorf("error getting the receipt from tx hash: %s with error: %w", hashTx, err)
	}

	var event eventBuyStorage

	for _, ev := range receipt.Logs {
		if ev.Address == s.grydContractAddress && len(ev.Topics) > 0 && ev.Topics[0] == s.Events.storageBoughtTopic {
			err = transaction.ParseEvent(&s.grydContractABI, "StorageBought", &event, *ev)
			if err != nil {
				return false, fmt.Errorf("no events found for the corresponding tx hash: %s with error: %w", hashTx, err)
			}
		}
	}

	testVal := fmt.Sprintf("Buyer Address: %s\nBuyerUsername: %s\nStorage Size: %d", event.Buyer.String(), event.UserName, event.Size.Int64())
	fmt.Println(testVal)
	return true, nil
}
