package server

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"github.com/gryd-database/platform-poc/configuration"
	"github.com/gryd-database/platform-poc/pkg/storage/dbMock"
	"github.com/gryd-database/platform-poc/pkg/storage/storageMock"
	"github.com/gryd-database/platform-poc/pkg/transaction/txMock"
	"github.com/sirupsen/logrus"
	"net/http"
	"testing"
)

type testServerOptions struct {
	config            *configuration.Config
	logger            *logrus.Logger
	router            *chi.Mux
	storageController *StorageController
	ethAddress        common.Address
	txServiceOpts     []txMock.Option
	dbServiceOpts     []dbMock.Option
	odbServiceOpts    []storageMock.Option
}

type extraOpts struct {
}

func newTestServer(t *testing.T, o testServerOptions) (*http.Client, *websocket.Conn, string) {
	t.Helper()
	//
	//transaction := txMock.New(o.txServiceOpts...)
	//
	//dbService := dbMock.New(o.dbServiceOpts...)
	//
	//storageService := storageMock.New(o.odbServiceOpts...)

	return nil, nil, ""
}
