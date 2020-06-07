package user

import (
	"fmt"

	"github.com/cynt4k/wygops/internal/event"
	"github.com/leandro-lugaresi/hub"
	"github.com/mitchellh/mapstructure"
)

type eventHandler func(ns *Service, ev hub.Message)

var handlerMap = map[string]eventHandler{
	event.UserCreated: userCreatedHandler,
}

func userCreatedHandler(ns *Service, ev hub.Message) {
	var msg event.UserCreatedEvent
	err := mapstructure.Decode(ev.Fields, &msg)
	if err != nil {
		ns.logger.Warn("error decode interface to struct")
		return
	}
	fmt.Println(msg)
}
