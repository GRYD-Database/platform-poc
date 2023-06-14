package server

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/gryd-database/platform-poc/cmd/server/controller"
	"github.com/gryd-database/platform-poc/configuration"
	"github.com/gryd-database/platform-poc/pkg/storage"
	"net/http"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"

	"github.com/gryd-database/platform-poc/pkg/cdb"
	"github.com/gryd-database/platform-poc/pkg/logger"
)

type Container struct {
	config *configuration.Config
	logger *logrus.Logger
	cdb    *pgxpool.Pool
	router *chi.Mux

	storageController *controller.StorageController
}

func Init() error {
	container, err := NewContainer()
	if err != nil {
		return fmt.Errorf("failed to initialize services: %w", err)
	}

	container.logger.Info("Container Initialized Successfully")

	container.storageController = controller.New(container.logger, storage.New(container.cdb, container.logger))

	container.router = chi.NewRouter()
	container.cors()
	container.routes()

	go func() {
		container.startServer()
	}()
	select {}
}

// NewContainer bootstrap important services
func NewContainer() (*Container, error) {
	confInstance, err := configuration.Init()
	if err != nil {
		return nil, fmt.Errorf("error bootstrapping config: %w", err)
	}

	loggerInstance, err := logger.Init(confInstance)
	if err != nil {
		return nil, fmt.Errorf("error bootstrapping logger: %w", err)
	}

	cdbInstance, err := cdb.Init(confInstance)

	return &Container{
		config: confInstance,
		logger: loggerInstance,
		cdb:    cdbInstance,
	}, nil
}

func (c *Container) routes() {
	c.router.Route("/storage", func(r chi.Router) {
		r.Post("/create", c.storageController.Create)
	})
}

func (c *Container) cors() {
	c.router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           int(12 * time.Hour),
	}))
}

func (c *Container) startServer() {
	c.logger.Info("Starting Server at:", c.config.Address)

	err := http.ListenAndServe(c.config.Address, c.router)
	if err != nil {
		c.logger.Info("error starting server at ", c.config.Address, " with error: ", err)
		panic(err)
	}
}
