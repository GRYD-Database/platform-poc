package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/gryd-database/platform-poc/configuration"
	"github.com/gryd-database/platform-poc/pkg/node"
	"github.com/gryd-database/platform-poc/pkg/odb"
	"github.com/gryd-database/platform-poc/pkg/pg"
	"github.com/gryd-database/platform-poc/pkg/storage"
	"github.com/gryd-database/platform-poc/pkg/transaction"
	"golang.org/x/sync/semaphore"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"

	"github.com/gryd-database/platform-poc/pkg/logger"
)

type Container struct {
	config            *configuration.Config
	logger            *logrus.Logger
	router            *chi.Mux
	pg                *pgxpool.Pool
	storageController *StorageController
	ethAddress        common.Address
	txService         *transaction.TxService
	odb               *odb.Database

	rpcClient     *rpc.Client
	grydSemaphore *semaphore.Weighted
}

func Init() error {
	container, err := NewContainer()
	if err != nil {
		return fmt.Errorf("failed to initialize services: %w", err)
	}

	container.logger.Info("Container Initialized Successfully")

	GRYDContractAddress, GRYDContractABI, err := setContracts(container.config.GRYDContract.Address, container.config.GRYDContract.ABI)
	if err != nil {
		container.logger.Error("failed to parse contract abi, ", err)
		return fmt.Errorf("err loading gryd contract: %w", err)
	}

	container.rpcClient, container.ethAddress, container.txService, err = node.InitChain(context.Background(), container.logger, container.config.ChainConfig.Endpoint, container.config.ChainConfig.PrivateKey)

	grydContract := storage.NewContract(container.txService, container.ethAddress, container.logger, GRYDContractAddress, GRYDContractABI)

	container.storageController = New(
		container.logger,
		storage.New(
			container.ethAddress,
			container.logger,
			container.pg, container.odb.Store, container.odb.Ledger),
		grydContract)

	container.router = chi.NewRouter()
	container.cors()
	container.routes()

	go func() {
		container.startServer()
	}()
	select {}
}

// NewContainer bootstrap important services
func NewContainer() (*Container, error) {
	confInstance, err := configuration.Init()
	if err != nil {
		return nil, fmt.Errorf("error bootstrapping config: %w", err)
	}

	loggerInstance, err := logger.Init(confInstance)
	if err != nil {
		return nil, fmt.Errorf("error bootstrapping logger: %w", err)
	}

	pgInstance, err := pg.InitPool(confInstance)
	if err != nil {
		loggerInstance.Error("failed to initialize pg instance: ", err)
		return nil, fmt.Errorf("error bootstrapping pg: %w", err)
	}

	odb, err := odb.NewDatabase(
		context.Background(),
		confInstance.IPFS.Address,
		confInstance.IPFS.RepoPath,
		confInstance.IPFS.IsLocal,
		confInstance.IPFS.CreateRepo,
		confInstance.IPFS.IsReplicated,
		loggerInstance,
	)
	if err != nil {
		loggerInstance.Error("failed to initialize odb/ipfs instance: ", err)
		return nil, fmt.Errorf("unable to bootstrap odb: %w", err)
	}

	return &Container{
		config: confInstance,
		logger: loggerInstance,
		pg:     pgInstance,
		odb:    odb,
	}, nil
}

func (c *Container) routes() {
	c.router.Route("/storage", func(r chi.Router) {
		c.grydAccessHandler()
		r.Post("/create", c.storageController.Create)
		r.Get("/get/{id}", c.storageController.GetRecordByID)
	})

	c.router.Route("/balance", func(r chi.Router) {
		c.grydAccessHandler()
		r.Get("/get", c.storageController.GetBalance)
	})
}

func (c *Container) cors() {
	c.router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           int(12 * time.Hour),
	}))
}

func (c *Container) startServer() {
	c.logger.Info("Starting Server at:", c.config.Address)

	err := http.ListenAndServe(c.config.Address, c.router)
	if err != nil {
		c.logger.Error("error starting server at ", c.config.Address, " with error: ", err)
		panic(err)
	}
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
