package server

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-chi/chi"
	"github.com/gryd-database/platform-poc/configuration"
	"github.com/gryd-database/platform-poc/pkg/storage/dbMock"
	"github.com/gryd-database/platform-poc/pkg/storage/grydContractMock"
	"github.com/gryd-database/platform-poc/pkg/storage/odbMock"
	"github.com/gryd-database/platform-poc/pkg/transaction/txMock"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testServerOptions struct {
	config                  *configuration.Config
	logger                  *logrus.Logger
	router                  *chi.Mux
	storageController       *StorageController
	ethAddress              common.Address
	txServiceOpts           []txMock.Option
	dbServiceOpts           []dbMock.Option
	odbServiceOpts          []odbMock.Option
	grydContractServiceOpts []grydContractMock.Option
}

type extraOpts struct {
}

func newTestServer(t *testing.T, o testServerOptions) *http.Client {
	t.Helper()

	transaction := txMock.New(o.txServiceOpts...)

	dbService := dbMock.New(o.dbServiceOpts...)

	storageService := odbMock.New(o.odbServiceOpts...)

	contractService := grydContractMock.New(o.grydContractServiceOpts...)

	storageController := New(logrus.New(), storageService, dbService, contractService)

	s := ContainerBootstrapper(nil, o.ethAddress, &transaction, nil, storageController)

	server := httptest.NewServer(s)

	return server.Client()
}
