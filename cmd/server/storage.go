package server

import (
	"context"
	"encoding/json"
	"errors"
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
		grydService:    grydContract,
	}
}

type StorageService interface {
	Create(ctx context.Context, voStorage *storage.VoStorage) (*storage.DTOStorage, error)
	AssignStorage() (string, error)
}

type GRYDContract interface {
	GetBalance(ctx context.Context) (*big.Int, error)
	VerifyEvent(ctx context.Context, hashTx string) (bool, error)
}

type StorageController struct {
	logger         *logrus.Logger
	storageService StorageService
	grydService    GRYDContract
}

func (c *StorageController) Create(w http.ResponseWriter, r *http.Request) {
	storageVo := storage.VoStorage{}

	err := json.NewDecoder(r.Body).Decode(&storageVo)
	if err != nil {
		c.logger.Info("invalid arguments in storageVo body")
		WriteJson(w, storage.DTOStorage{}, http.StatusBadRequest)
		return
	}

	reInput := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	if !(reInput.MatchString(storageVo.Wallet)) {
		c.logger.Info("invalid wallet address:" + storageVo.Wallet)

		WriteJson(w, storage.DTOStorage{}, http.StatusBadRequest)
		return
	}

	reInput = regexp.MustCompile("^0x([A-Fa-f0-9]{64})$")
	if !(reInput.MatchString(storageVo.TxHash)) {
		c.logger.Info("invalid tx hash:" + storageVo.TxHash)

		WriteJson(w, storage.DTOStorage{}, http.StatusBadRequest)
		return
	}

	verified, err := c.grydService.VerifyEvent(r.Context(), storageVo.TxHash)
	if err != nil {
		if errors.Is(err, transaction.ErrEventNotFound) {
			c.logger.Info("event not found for tx hash:" + storageVo.TxHash)

			WriteJson(w, storage.DTOStorage{}, http.StatusNotFound)
			return
		}

		if errors.Is(err, transaction.ErrNoTopic) {
			c.logger.Info("topic not found for tx hash:" + storageVo.TxHash)

			WriteJson(w, storage.DTOStorage{}, http.StatusNotFound)
			return
		}

		if errors.Is(err, storage.ErrUnprocessableEvent) {
			c.logger.Info("tx receipt or event does not exist for hash:" + storageVo.TxHash)

			WriteJson(w, storage.DTOStorage{}, http.StatusNotFound)
			return
		}

		c.logger.Error("internal server error: ", err)
		WriteJson(w, storage.DTOStorage{}, http.StatusInternalServerError)
		return
	}

	if !verified {
		c.logger.Info("cannot verify event for tx: ", storageVo.TxHash)

		WriteJson(w, storage.DTOStorage{}, http.StatusBadRequest)
		return
	}

	dsn, _ := c.storageService.AssignStorage()
	println(dsn)

	resp, err := c.storageService.Create(r.Context(), &storageVo)
	if err != nil {
		c.logger.Error("internal server error: ", err)

		WriteJson(w, storage.DTOStorage{}, http.StatusInternalServerError)
		return
	}

	WriteJson(w, resp, http.StatusOK)
}

func (c *StorageController) GetBalance(w http.ResponseWriter, r *http.Request) {
	balance, err := c.grydService.GetBalance(r.Context())
	if err != nil {
		c.logger.Error("internal server error: ", err)
		WriteJson(w, storage.VoStorage{}, http.StatusInternalServerError)
		return
	}

	WriteJson(w, balance, http.StatusOK)
}
