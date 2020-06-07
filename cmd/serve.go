package cmd

import (
	"fmt"

	"github.com/cynt4k/wygops/internal/repository"
	service "github.com/cynt4k/wygops/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/leandro-lugaresi/hub"
	"go.uber.org/zap"
)

// HTTPServer : HttpServer struct to initialize
type HTTPServer struct {
	Logger *zap.Logger
	SS     *service.Services
	Router *gin.Engine
	Hub    *hub.Hub
	Repo   repository.Repository
}

// ServeServer : Execute the server
func ServeServer() error {
	// New Message bus
	hub := hub.New()
	logger := getLogger()

	logger.Info("connecting to database...")
	db, err := c.getDatabase()

	if err != nil {
		logger.Fatal("error while connecting to the database", zap.Error(err))
	}
	logger.Info("connection to database established")

	logger.Info("setting up repository...")
	repo, err := repository.NewGormRepository(db, hub)

	if err != nil {
		logger.Fatal("error while setting up repository", zap.Error(err))
	}
	logger.Info("repository initialized")

	logger.Info("sync the repo..")
	synced, err := repo.Sync()

	if err != nil {
		logger.Fatal("error while syncing the repo", zap.Error(err))
	}

	if synced {
		logger.Info("repository is synced")
	} else {
		logger.Info("repository was not synced")
	}

	server, err := newHttpServer(hub, db, repo, logger)

	if err != nil {
		logger.Fatal("error while creating server", zap.Error(err))
	}
	go func() {
		if err := server.Start(fmt.Sprintf(":%d", c.API.Port)); err != nil {
			logger.Info("shutting down the api")
		}
	}()

	return err
}

// Start : Start the HTTPServer
func (s *HTTPServer) Start(address string) error {
	return s.Router.Run(address)
}
