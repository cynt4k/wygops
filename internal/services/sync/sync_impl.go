package sync

import (
	"fmt"

	"github.com/cynt4k/wygops/internal/repository"
	"github.com/cynt4k/wygops/internal/services/ldap"
	"github.com/leandro-lugaresi/hub"
	"go.uber.org/zap"
)

// Service : Sync service struct
type Service struct {
	hub        *hub.Hub
	repo       repository.Repository
	ldap       *ldap.LDAP
	logger     *zap.Logger
	sourceType string
}

// NewService : Create a new sync service
func NewService(repo repository.Repository, hub *hub.Hub, source interface{}, logger *zap.Logger) (Sync, error) {
	service := &Service{
		hub:  hub,
		repo: repo,
	}

	switch sourceType := source.(type) {
	case *ldap.LDAP:
		service.ldap = sourceType
		service.sourceType = "ldap"
		break
	default:
		return nil, fmt.Errorf("no valid source type added")
	}

	return service, nil
}

// Start : Start the sync service
func (s *Service) Start() error {
	return nil
}

// Stop : Stop the sync service
func (s *Service) Stop() error {
	return nil
}
