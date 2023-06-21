package node

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	signer2 "github.com/gryd-database/platform-poc/pkg/signer"
	"github.com/gryd-database/platform-poc/pkg/transaction"
	"github.com/sirupsen/logrus"
)

func InitChain(
	ctx context.Context,
	logger *logrus.Logger,
	endpoint string,
	hexKey string) (common.Address, *transaction.TxService, error) {

	rpcClient, err := rpc.DialContext(ctx, endpoint)
	if err != nil {
		logger.Error("unable to dial eth client: ", err)
		return common.Address{}, nil, fmt.Errorf("unable to dial eth client: %w", err)
	}

	var versionString string
	err = rpcClient.CallContext(ctx, &versionString, "web3_clientVersion")
	if err != nil {
		logger.Info("could not connect to backend, requires a working blockchain note, please check your endpoint: endpoint=", endpoint, "\n error: ", err)
		return common.Address{}, nil, fmt.Errorf("could not connect to backend, requires a working blockchain note, please check your endpoint: err: %w", err)
	}

	backend := transaction.Backend(ethclient.NewClient(rpcClient))

	chainID, err := backend.ChainID(ctx)
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("could not get chain ID")
	}

	signer, err := signer2.New(hexKey)
	if err != nil {
		return common.Address{}, nil, err
	}

	ethAddress, err := signer.EthereumAddress()
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("error getting eth address: %w", err)
	}

	txService, err := transaction.NewTxService(logger, backend, signer, chainID, ethAddress)
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("error bootstrapping transaction service: %w", err)
	}

	return ethAddress, txService, nil
}
