package server

import (
	"context"
	"encoding/csv"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/gryd-database/platform-poc/pkg/storage"
	"github.com/gryd-database/platform-poc/pkg/transaction"
	"github.com/sirupsen/logrus"
	"math/big"
	"net/http"
	"regexp"
)

func New(logger *logrus.Logger, service *storage.Storage, grydContract *storage.Contract) *StorageController {
	return &StorageController{
		logger:         logger,
		storageService: service,
		odbService:     service,
		grydService:    grydContract,
	}
}

type StorageService interface {
	Create(ctx context.Context, voStorage *storage.VoStorage) (*storage.DTOStorage, error)
}

type OrbitService interface {
	AddRecord(ctx context.Context, storage *[]storage.InputData) error
	Ledger(ctx context.Context, wallet, datasetKey string) error
	GetWalletByDatasetKey(ctx context.Context, key string) (*storage.Ledger, error)
	GetRecordByID(ctx context.Context, id string) (*storage.InputData, error)
}

type GRYDContract interface {
	GetBalance(ctx context.Context) (*big.Int, error)
	VerifyEvent(ctx context.Context, hashTx string) (*storage.EventBuyStorage, error)
}

type StorageController struct {
	logger         *logrus.Logger
	storageService StorageService
	grydService    GRYDContract
	odbService     OrbitService
}

func (c *StorageController) Create(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("file")
	if err != nil {
		c.logger.Info("unable to parse form data: ", err)

		WriteJson(w, "unable to parse form data", http.StatusInternalServerError)
		return
	}

	var inputDataObject []storage.InputData

	reader := csv.NewReader(file)
	record, err := reader.ReadAll()
	if err != nil {
		c.logger.Info("unable to parse form data: ", err)

		WriteJson(w, "unable to parse form data", http.StatusInternalServerError)
		return
	}

	datasetKey := uuid.NewString()

	storageVo := storage.VoStorage{
		Wallet:     r.FormValue("wallet"),
		TxHash:     r.FormValue("txHash"),
		DatasetKey: datasetKey,
	}

	for _, line := range record {
		inputData := storage.InputData{
			ID:         uuid.NewString(),
			Date:       line[1],
			DataType:   line[2],
			Data:       line[3],
			Dataset:    line[0],
			DatasetKey: datasetKey,
		}
		inputDataObject = append(inputDataObject, inputData)
	}

	reInput := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	if !(reInput.MatchString(storageVo.Wallet)) {
		c.logger.Info("invalid wallet address:" + storageVo.Wallet)

		WriteJson(w, "invalid wallet address", http.StatusBadRequest)
		return
	}

	reInput = regexp.MustCompile("^0x([A-Fa-f0-9]{64})$")
	if !(reInput.MatchString(storageVo.TxHash)) {
		c.logger.Info("invalid tx hash:" + storageVo.TxHash)

		WriteJson(w, "invalid tx hash", http.StatusBadRequest)
		return
	}

	event, err := c.grydService.VerifyEvent(r.Context(), storageVo.TxHash)
	if err != nil {
		if errors.Is(err, transaction.ErrEventNotFound) {
			c.logger.Info("event not found for tx hash:" + storageVo.TxHash)

			WriteJson(w, "event not found", http.StatusNotFound)
			return
		}

		if errors.Is(err, transaction.ErrNoTopic) {
			c.logger.Info("topic not found for tx hash:" + storageVo.TxHash)

			WriteJson(w, "event cannot be processed", http.StatusNotFound)
			return
		}

		if errors.Is(err, storage.ErrUnprocessableEvent) {
			c.logger.Info("tx receipt or event does not exist for hash:" + storageVo.TxHash)

			WriteJson(w, "tx receipt or event does not exist for hash", http.StatusNotFound)
			return
		}

		c.logger.Error("internal server error: ", err)
		WriteJson(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if event.Buyer != common.HexToAddress(storageVo.Wallet) {
		c.logger.Info("cannot verify event for tx: ", storageVo.TxHash)

		WriteJson(w, "cannot verify event for tx", http.StatusBadRequest)
		return
	}

	err = c.odbService.AddRecord(r.Context(), &inputDataObject)
	if err != nil {
		c.logger.Error("internal server error: ", err)

		WriteJson(w, "internal server error", http.StatusInternalServerError)
		return
	}

	err = c.odbService.Ledger(r.Context(), storageVo.Wallet, datasetKey)
	if err != nil {
		c.logger.Error("internal server error: ", err)

		WriteJson(w, "internal server error", http.StatusInternalServerError)
		return
	}

	resp, err := c.storageService.Create(r.Context(), &storageVo)
	if err != nil {
		c.logger.Error("internal server error: ", err)

		WriteJson(w, "internal server error", http.StatusInternalServerError)
		return
	}

	WriteJson(w, resp, http.StatusOK)
}

func (c *StorageController) GetBalance(w http.ResponseWriter, r *http.Request) {
	balance, err := c.grydService.GetBalance(r.Context())
	if err != nil {
		c.logger.Error("internal server error: ", err)
		WriteJson(w, "internal server error", http.StatusInternalServerError)
		return
	}

	WriteJson(w, balance, http.StatusOK)
}

func (c *StorageController) GetRecordByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if len(id) == 0 {
		c.logger.Error("id is missing in path params")
		WriteJson(w, "id is missing in path params", http.StatusBadRequest)
		return
	}

	record, err := c.odbService.GetRecordByID(r.Context(), id)
	if err != nil {
		c.logger.Error("internal server error: ", err)
		WriteJson(w, "internal server error", http.StatusInternalServerError)
		return
	}

	WriteJson(w, record, http.StatusOK)
}

// 1553bd83-f1cf-47f4-9a22-1f7251f4dfcd
// 3b3ae135-a364-462d-8f8f-65ccdb730600
