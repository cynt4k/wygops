package sync

import (
	"github.com/cynt4k/wygops/internal/event"
	"github.com/leandro-lugaresi/hub"
	"github.com/mitchellh/mapstructure"
)

type eventHandler func(s *Service, ev hub.Message)

var handlerMap = map[string]eventHandler{
	event.UserCreated: userCreatedHandler,
}

func userCreatedHandler(s *Service, ev hub.Message) {
	var msg event.UserCreatedEvent
	err := mapstructure.Decode(ev.Fields, &msg)
	if err != nil {
		s.logger.Warn("error decode interface to struct")
		return
	}

	// TODO: trigger manual sync
}
