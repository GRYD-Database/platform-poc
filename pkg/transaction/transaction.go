package transaction

import (
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/gryd-database/platform-poc/pkg/signer"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"io"
	"math/big"
	"sync"
)

var (
	// ErrTransactionReverted denotes that the sent transaction has been
	// reverted.
	ErrTransactionReverted = errors.New("transaction reverted")
	ErrUnknownTransaction  = errors.New("unknown transaction")
	ErrAlreadyImported     = errors.New("already imported")
)

// TxRequest describes a request for a transaction that can be executed.
type TxRequest struct {
	To                   *common.Address // recipient of the transaction
	Data                 []byte          // transaction data
	GasPrice             *big.Int        // gas price or nil if suggested gas price should be used
	GasLimit             uint64          // gas limit or 0 if it should be estimated
	MinEstimatedGasLimit uint64          // minimum gas limit to use if the gas limit was estimated; it will not apply when this value is 0 or when GasLimit is not 0
	GasFeeCap            *big.Int        // adds a cap to maximum fee user is willing to pay
	Value                *big.Int        // amount of wei to send
	Description          string          // optional description
}

type StoredTransaction struct {
	To          *common.Address // recipient of the transaction
	Data        []byte          // transaction data
	GasPrice    *big.Int        // used gas price
	GasLimit    uint64          // used gas limit
	GasTipBoost int             // adds a tip for the miner for prioritizing transaction
	GasTipCap   *big.Int        // adds a cap to the tip
	GasFeeCap   *big.Int        // adds a cap to maximum fee user is willing to pay
	Value       *big.Int        // amount of wei to send
	Nonce       uint64          // used nonce
	Created     int64           // creation timestamp
	Description string          // description
}

// Service is the service to send transactions. It takes care of gas price, gas
// limit and nonce management.
type Service interface {
	io.Closer
	// Send creates a transaction based on the request (with gasprice increased by provided percentage) and sends it.
	Send(ctx context.Context, request *TxRequest, tipCapBoostPercent int) (txHash common.Hash, err error)
	// Call simulate a transaction based on the request.
	Call(ctx context.Context, request *TxRequest) (result []byte, err error)
	// WaitForReceipt waits until either the transaction with the given hash has been mined or the context is cancelled.
	// This is only valid for transaction sent by this service.
	WaitForReceipt(ctx context.Context, txHash common.Hash) (receipt *types.Receipt, err error)
	// WatchSentTransaction start watching the given transaction.
	// This wraps the monitors watch function by loading the correct nonce from the store.
	// This is only valid for transaction sent by this service.
	WatchSentTransaction(txHash common.Hash) (<-chan types.Receipt, <-chan error, error)
	// StoredTransaction retrieves the stored information for the transaction
	StoredTransaction(txHash common.Hash) (*StoredTransaction, error)
	// PendingTransactions retrieves the list of all pending transaction hashes
	PendingTransactions() ([]common.Hash, error)
	// ResendTransaction resends a previously sent transaction
	// This operation can be useful if for some reason the transaction vanished from the eth networks pending pool
	ResendTransaction(ctx context.Context, txHash common.Hash) error
	// CancelTransaction cancels a previously sent transaction by double-spending its nonce with zero-transfer one
	CancelTransaction(ctx context.Context, originalTxHash common.Hash) (common.Hash, error)
	// TransactionFee retrieves the transaction fee
	TransactionFee(ctx context.Context, txHash common.Hash) (*big.Int, error)
}

type TxService struct {
	wg     sync.WaitGroup
	lock   sync.Mutex
	ctx    context.Context
	cancel context.CancelFunc

	logger  *logrus.Logger
	backend Backend
	signer  signer.Signer
	sender  common.Address
	chainID *big.Int
}

func (t *TxService) Close() error {
	//TODO implement me
	panic("implement me")
}

func (t *TxService) Send(ctx context.Context, request *TxRequest, tipCapBoostPercent int) (txHash common.Hash, err error) {
	//TODO implement me
	panic("implement me")
}

func (t *TxService) Call(ctx context.Context, request *TxRequest) (result []byte, err error) {
	msg := ethereum.CallMsg{
		From:     t.sender,
		To:       request.To,
		Data:     request.Data,
		GasPrice: request.GasPrice,
		Gas:      request.GasLimit,
		Value:    request.Value,
	}
	data, err := t.backend.CallContract(ctx, msg, nil)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (t *TxService) WaitForReceipt(ctx context.Context, txHash common.Hash) (receipt *types.Receipt, err error) {
	//TODO implement me
	panic("implement me")
}

func (t *TxService) WatchSentTransaction(txHash common.Hash) (<-chan types.Receipt, <-chan error, error) {
	//TODO implement me
	panic("implement me")
}

func (t *TxService) StoredTransaction(txHash common.Hash) (*StoredTransaction, error) {
	//TODO implement me
	panic("implement me")
}

func (t *TxService) PendingTransactions() ([]common.Hash, error) {
	//TODO implement me
	panic("implement me")
}

func (t *TxService) ResendTransaction(ctx context.Context, txHash common.Hash) error {
	//TODO implement me
	panic("implement me")
}

func (t *TxService) CancelTransaction(ctx context.Context, originalTxHash common.Hash) (common.Hash, error) {
	//TODO implement me
	panic("implement me")
}

func (t *TxService) TransactionFee(ctx context.Context, txHash common.Hash) (*big.Int, error) {
	//TODO implement me
	panic("implement me")
}

func NewTxService(logger *logrus.Logger, backend Backend, signer signer.Signer, chainID *big.Int, address common.Address) (*TxService, error) {
	ctx, cancel := context.WithCancel(context.Background())

	return &TxService{
		wg:      sync.WaitGroup{},
		lock:    sync.Mutex{},
		ctx:     ctx,
		cancel:  cancel,
		logger:  logger,
		backend: backend,
		signer:  signer,
		sender:  address,
		chainID: chainID,
	}, nil
}
