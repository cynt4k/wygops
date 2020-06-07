package user

import (
	"github.com/cynt4k/wygops/internal/repository"
	"github.com/leandro-lugaresi/hub"
	"go.uber.org/zap"
)

// Service : Service struct
type Service struct {
	repo   repository.Repository
	hub    *hub.Hub
	logger *zap.Logger
}

// NewService : Create a new user service
func NewService(repo repository.Repository, hub *hub.Hub, logger *zap.Logger) *Service {
	service := &Service{
		repo:   repo,
		hub:    hub,
		logger: logger,
	}
	go func() {
		topics := make([]string, 0, len(handlerMap))
		for k := range handlerMap {
			topics = append(topics, k)
		}
		for msg := range hub.Subscribe(200, topics...).Receiver {
			h, ok := handlerMap[msg.Topic()]
			if ok {
				go h(service, msg)
			}
		}
	}()
	return service
}
