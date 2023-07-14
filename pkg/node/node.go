package node

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/gryd-database/platform-poc/pkg/signer"
	"github.com/gryd-database/platform-poc/pkg/transaction"
	"github.com/sirupsen/logrus"
)

func InitChain(
	ctx context.Context,
	logger *logrus.Logger,
	endpoint string,
	hexKey string) (*rpc.Client, common.Address, *transaction.Service, error) {

	rpcClient, err := rpc.DialContext(ctx, endpoint)
	if err != nil {
		logger.Error("unable to dial eth client: ", err)
		return nil, common.Address{}, nil, fmt.Errorf("unable to dial eth client: %w", err)
	}

	var versionString string

	err = rpcClient.CallContext(ctx, &versionString, "web3_clientVersion")
	if err != nil {
		logger.Info("could not connect to backend, requires a working blockchain note, please check your endpoint: endpoint=", endpoint, "\n error: ", err)
		return nil, common.Address{}, nil, fmt.Errorf("could not connect to backend, requires a working blockchain note, please check your endpoint: err: %w", err)
	}

	backend := transaction.NewBackend(ethclient.NewClient(rpcClient))

	chainID, err := backend.ChainID(ctx)
	if err != nil {
		return nil, common.Address{}, nil, fmt.Errorf("could not get chain ID")
	}

	signer, err := signer.New(hexKey)
	if err != nil {
		return nil, common.Address{}, nil, err
	}

	ethAddress := signer.EthereumAddress()

	txService, err := transaction.NewTxService(rpcClient, logger, *backend, signer, chainID, ethAddress)
	if err != nil {
		return nil, common.Address{}, nil, fmt.Errorf("error bootstrapping transaction service: %w", err)
	}

	return rpcClient, ethAddress, &txService, nil
}
