package server

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-chi/chi"
	"github.com/gryd-database/platform-poc/configuration"
	"github.com/gryd-database/platform-poc/pkg/storage"
	"github.com/gryd-database/platform-poc/pkg/transaction/txMock"
	"github.com/sirupsen/logrus"
	"testing"
)

type testServerOptions struct {
	config                  *configuration.Config
	logger                  *logrus.Logger
	router                  *chi.Mux
	storageController       *StorageController
	ethAddress              common.Address
	txServiceOpts           []txMock.Option
	dbServiceOpts           storage.DBService
	odbServiceOpts          storage.OrbitService
	grydContractServiceOpts storage.GRYDContract
}

func newTestServer(t *testing.T, o testServerOptions) *Container {
	t.Helper()

	transaction := txMock.New(o.txServiceOpts...)

	dbService := o.dbServiceOpts

	storageService := o.odbServiceOpts

	contractService := o.grydContractServiceOpts

	config, err := configuration.Init()
	if err != nil {
		t.Fatal(err)
	}

	storageController := New(logrus.New(), storageService, dbService, contractService)

	s := ContainerBootstrapper(nil, o.ethAddress, &transaction, &BootedServices{config: config}, storageController)

	s.cors()
	s.routes()

	return s
}
