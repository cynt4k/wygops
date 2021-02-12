package sync

import (
	"fmt"
	"time"

	"github.com/cynt4k/wygops/cmd/config"
	"github.com/cynt4k/wygops/internal/repository"
	"github.com/cynt4k/wygops/internal/services/ldap"
	"github.com/leandro-lugaresi/hub"
	"go.uber.org/zap"
)

// Service : Sync service struct
type Service struct {
	hub           *hub.Hub
	repo          repository.Repository
	ldap          ldap.LDAP
	logger        *zap.Logger
	duration      time.Duration
	syncInterrupt chan bool
	sourceType    string
}

// NewLDAPService : Create a new sync service
func NewLDAPService(
	repo repository.Repository,
	hub *hub.Hub,
	source ldap.LDAP,
	logger *zap.Logger,
	config *config.Config,
) (Sync, error) {
	return newService(repo, hub, source, logger, config)
}

// newService : Create a new sync service
func newService(
	repo repository.Repository,
	hub *hub.Hub,
	source interface{},
	logger *zap.Logger,
	config *config.Config,
) (Sync, error) {
	service := &Service{
		hub:    hub,
		repo:   repo,
		logger: logger,
	}

	duration, err := time.ParseDuration(config.General.Sync.Interval)

	if err != nil {
		return nil, err
	}

	service.duration = duration

	switch sourceType := source.(type) {
	case ldap.LDAP:
		service.ldap = sourceType
		service.sourceType = "ldap"
	default:
		return nil, fmt.Errorf("no valid source type added")
	}

	err = service.Start()

	if err != nil {
		return nil, err
	}

	service.initEventhandler()

	return service, nil
}

// Start : Start the sync service
func (s *Service) Start() error {
	switch s.sourceType {
	case "ldap":
		return s.startLdap()
	default:
		return fmt.Errorf("unknown source type - check the initialization")
	}
}

// Stop : Stop the sync service
func (s *Service) Stop() error {
	return nil
}

func (s *Service) initEventhandler() {
	go func() {
		const capSize = 200
		topics := make([]string, 0, len(handlerMap))
		for k := range handlerMap {
			topics = append(topics, k)
		}
		for msg := range s.hub.Subscribe(capSize, topics...).Receiver {
			h, ok := handlerMap[msg.Topic()]
			if ok {
				go h(s, msg)
			}
		}
	}()
}
