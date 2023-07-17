package server

import (
	"bytes"
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-chi/chi"
	"github.com/gryd-database/platform-poc/pkg/storage"
	"github.com/gryd-database/platform-poc/pkg/storage/grydContractMock"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
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

		res := httptest.NewRecorder()

		// router is of type http.Handler
		router.ServeHTTP(res, req)

		assertStatusCode(t, res.Code, http.StatusBadRequest)
	})
}

func assertStatusCode(t testing.TB, got, want int) {
	t.Helper()

	if got != want {
		t.Errorf("handler returned wrong status code: got %v want %v", got, want)
	}
}
