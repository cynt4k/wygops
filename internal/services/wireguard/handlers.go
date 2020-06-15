package wireguard

import (
	"fmt"
	"net"

	"github.com/cynt4k/wygops/internal/event"
	"github.com/leandro-lugaresi/hub"
	"github.com/mitchellh/mapstructure"

	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

type eventHandler func(w *Service, ev hub.Message)

var handlerMap = map[string]eventHandler{
	event.DeviceCreated: deviceCreatedHandler,
}

func deviceCreatedHandler(w *Service, ev hub.Message) {
	var msg event.DeviceCreatedEvent
	err := mapstructure.Decode(ev.Fields, &msg)
	if err != nil {
		w.logger.Warn("error decode interface to struct")
		return
	}

	device, err := w.repo.GetDevice(msg.DeviceID)

	if err != nil {
		w.logger.Warn(fmt.Sprintf("Device not found with id %v", msg.DeviceID))
		return
	}

	publicKey, err := wgtypes.ParseKey(device.PublicKey)

	if err != nil {
		w.logger.Warn(fmt.Sprintf("Public key for device %v invalid", msg.DeviceID))
		return
	}

	peer := Peer{
		PublicKey:   publicKey,
		IPV4Address: net.ParseIP(device.IPv4Address),
		IPV6Address: net.ParseIP(device.IPv6Address),
	}
	w.addPeer(&peer)
}

func deviceDeletedHandler(w *Service, ev hub.Message) {
	var msg event.DeviceDeletedEvent
	err := mapstructure.Decode(ev.Fields, &msg)
	if err != nil {
		w.logger.Warn("error decode interface to struct")
		return
	}

	publicKey, err := wgtypes.ParseKey(msg.PublicKey)

	if err != nil {
		w.logger.Warn(fmt.Sprintf("Public key for device %v is invalid", msg.DeviceID))
		return
	}

	peer := Peer{
		PublicKey: publicKey,
	}
	w.deletePeer(&peer)
}

func deviceUpdatedHandler(w *Service, ev hub.Message) {
	var msg event.DeviceUpdatedEvent
	err := mapstructure.Decode(ev.Fields, &msg)
	if err != nil {
		w.logger.Warn("error decode interface to struct")
		return
	}

	device, err := w.repo.GetDevice(msg.DeviceID)

	if err != nil {
		w.logger.Warn(fmt.Sprintf("Device not found with id %v", msg.DeviceID))
		return
	}

	publicKey, err := wgtypes.ParseKey(device.PublicKey)

	if err != nil {
		w.logger.Warn(fmt.Sprintf("Public key for device %v is invalid", msg.DeviceID))
		return
	}

	peer := Peer{
		PublicKey:   publicKey,
		IPV4Address: net.ParseIP(device.IPv4Address),
		IPV6Address: net.ParseIP(device.IPv6Address),
	}
	w.updatePeer(&peer)
}
