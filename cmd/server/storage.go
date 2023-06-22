package server

import (
	"context"
	"encoding/json"
	"github.com/gryd-database/platform-poc/pkg/storage"
	"github.com/sirupsen/logrus"
	"math/big"
	"net/http"
)

func New(logger *logrus.Logger, service *storage.Storage) *StorageController {
	return &StorageController{
		logger:         logger,
		storageService: service,
	}
}

type StorageService interface {
	Create(ctx context.Context, voStorage *storage.VoStorage) (*storage.DTOStorage, error)
	GetBalance(ctx context.Context) (*big.Int, error)
}

type StorageController struct {
	logger         *logrus.Logger
	storageService *storage.Storage
}

func (c *StorageController) Create(w http.ResponseWriter, r *http.Request) {
	storageVo := storage.VoStorage{}

	err := json.NewDecoder(r.Body).Decode(&storageVo)
	if err != nil {
		c.logger.Info("invalid arguments in storageVo body")
		WriteJson(w, storage.VoStorage{}, http.StatusBadRequest)
		return
	}

	resp, err := c.storageService.Create(r.Context(), &storageVo)
	if err != nil {
		c.logger.Error("internal server error: ", err)
		WriteJson(w, storage.VoStorage{}, http.StatusInternalServerError)
		return
	}

	WriteJson(w, resp, http.StatusOK)
}

func (c *StorageController) GetBalance(w http.ResponseWriter, r *http.Request) {
	balance, err := c.storageService.GetBalance(r.Context())
	if err != nil {
		c.logger.Error("internal server error: ", err)
		WriteJson(w, storage.VoStorage{}, http.StatusInternalServerError)
		return
	}

	WriteJson(w, balance, http.StatusOK)
}