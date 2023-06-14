package controller

import (
	"github.com/gryd-database/platform-poc/pkg/storage"
	"github.com/sirupsen/logrus"
	"net/http"
)

func New(logger *logrus.Logger, service *storage.Storage) *StorageController {
	return &StorageController{
		logger:         logger,
		storageService: service,
	}
}

type StorageService interface {
	Create() error
}

type StorageController struct {
	logger         *logrus.Logger
	storageService *storage.Storage
}

func (c *StorageController) Create(w http.ResponseWriter, r *http.Request) {
	WriteJson(w, "success", http.StatusOK)
}
