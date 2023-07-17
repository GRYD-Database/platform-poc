package server

import (
	"bytes"
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/gryd-database/platform-poc/pkg/storage"
	"github.com/gryd-database/platform-poc/pkg/storage/dbMock"
	"github.com/gryd-database/platform-poc/pkg/storage/grydContractMock"
	"github.com/gryd-database/platform-poc/pkg/storage/odbMock"
	"github.com/gryd-database/platform-poc/pkg/transaction"
	"github.com/gryd-database/platform-poc/pkg/transaction/txMock"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestCreateStorage(t *testing.T) {
	t.Parallel()

	txHash := common.HexToHash("0xcb0caeff88b8bda3656396b19b808cd8b35c0054e96553852441ea2c3f5f4d26")
	address := "0xD07708adfbE343297E2ABfb31534dD3d78fg452a"
	createStorage := func() string {
		return fmt.Sprintf("/storage/create")
	}

	t.Run("ok", func(t *testing.T) {
		t.Parallel()

		event := &storage.EventInsertDataSuccess{
			User:      common.HexToAddress(address),
			QueryType: "create",
		}

		contract := grydContractMock.New(
			grydContractMock.WithVerifyEvent(func(ctx context.Context, hashTx string) (*storage.EventInsertDataSuccess, error) {
				return event, nil
			}))
		
		dbService := dbMock.New(
			dbMock.WithCreate(func(ctx context.Context, voStorage *storage.VoStorage) (*storage.DTOStorage, error) {
				return &storage.DTOStorage{
					ID:         uuid.UUID{},
					Wallet:     address,
					TxHash:     txHash.String(),
					CreatedAt:  time.Now(),
					DatasetKey: uuid.NewString(),
				}, nil
			}))
		
		odbService := odbMock.New(
			odbMock.WithAddRecord(func(ctx context.Context, storage *[]storage.InputData) error {
				
			})
		filePath := "../../sampleData.csv"
		fieldName := "file"
		body := new(bytes.Buffer)

		mw := multipart.NewWriter(body)

		file, err := os.Open(filePath)
		if err != nil {
			t.Fatal(err)
		}

		w, err := mw.CreateFormFile(fieldName, filePath)
		if err != nil {
			t.Fatal(err)
		}

		if _, err := io.Copy(w, file); err != nil {
			t.Fatal(err)
		}

		wallet, err := mw.CreateFormField("wallet")
		if err != nil {
			t.Fatal(err)
		}

		wallet.Write([]byte(address))

		tx, err := mw.CreateFormField("txHash")
		if err != nil {
			t.Fatal(err)
		}

		tx.Write(txHash.Bytes())

		// close the writer before making the request
		mw.Close()

		router := chi.NewRouter()
		req := httptest.NewRequest(http.MethodPost, createStorage(), body)

		req.Header.Add("Content-Type", mw.FormDataContentType())

		testServer := newTestServer()
	})
}

func assertStatusCode(t testing.TB, got, want int) {
	t.Helper()

	if got != want {
		t.Errorf("handler returned wrong status code: got %v want %v", got, want)
	}
}
