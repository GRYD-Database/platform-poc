package server

import (
	"bytes"
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/gryd-database/platform-poc/pkg/storage"
	"github.com/gryd-database/platform-poc/pkg/storage/dbMock"
	"github.com/gryd-database/platform-poc/pkg/storage/grydContractMock"
	"github.com/gryd-database/platform-poc/pkg/storage/odbMock"
	"github.com/magiconair/properties/assert"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func TestCreateStorage(t *testing.T) {
	t.Parallel()

	txHash := common.HexToHash("0xcb0caeff88b8bda3656396b19b808cd8b35c0054e96553852441ea2c3f5f4d26")
	address := "0xD07708ad91fbE34329507E2adABfb31534dD3efd"
	createStorage := func() string {
		return fmt.Sprintf("/storage/create")
	}

	t.Run("ok", func(t *testing.T) {
		t.Parallel()

		event := &storage.EventInsertDataSuccess{
			User:      common.HexToAddress(address),
			QueryType: "create",
		}

		id, _ := uuid.NewUUID()
		datasetKey := uuid.NewString()

		contract := grydContractMock.New(
			grydContractMock.WithVerifyEvent(func(ctx context.Context, hashTx string) (*storage.EventInsertDataSuccess, error) {
				return event, nil
			}))

		dbService := dbMock.New(
			dbMock.WithCreate(func(ctx context.Context, voStorage *storage.VoStorage) (*storage.DTOStorage, error) {
				return &storage.DTOStorage{
					ID:         id,
					Wallet:     address,
					TxHash:     txHash.String(),
					CreatedAt:  time.Now(),
					DatasetKey: datasetKey,
				}, nil
			}))

		odbService := odbMock.New(
			odbMock.WithAddRecord(func(ctx context.Context, storage *[]storage.InputData) error {
				return nil
			}),
			odbMock.WithLedger(func(ctx context.Context, wallet, datasetKey string) error {
				return nil
			}))

		testServer := newTestServer(t, testServerOptions{odbServiceOpts: odbService, dbServiceOpts: dbService, grydContractServiceOpts: contract})

		//prepare the reader instances to encode
		v := map[string]io.Reader{
			"file":   mustOpen("../../sampleData.csv"),
			"wallet": strings.NewReader(address),
			"txHash": strings.NewReader(txHash.String()),
		}

		req, err := Upload(v, createStorage())
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		testServer.router.ServeHTTP(rr, req)
		assert.Equal(t, rr.Result().StatusCode, http.StatusOK)
	})
}

func Upload(values map[string]io.Reader, url string) (req *http.Request, err error) {
	// Prepare a form that you will submit to that URL.
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for key, r := range values {
		var fw io.Writer
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}

		if x, ok := r.(*os.File); ok {
			if fw, err = w.CreateFormFile(key, x.Name()); err != nil {
				return
			}
		} else {
			// Add other fields
			if fw, err = w.CreateFormField(key); err != nil {
				return
			}
		}
		if _, err = io.Copy(fw, r); err != nil {
			return nil, err
		}

	}
	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	w.Close()

	// Now that you have a form, you can submit it to your handler.
	req, err = http.NewRequest("POST", url, &b)
	if err != nil {
		return
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Add("Content-Type", w.FormDataContentType())
	fmt.Println(w.FormDataContentType())
	return req, nil
}

func mustOpen(f string) *os.File {
	r, err := os.Open(f)
	if err != nil {
		panic(err)
	}
	return r
}
