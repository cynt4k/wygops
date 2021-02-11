package wireguard

import (
	"fmt"
	"net"

	"github.com/cynt4k/wygops/internal/event"
	"github.com/cynt4k/wygops/internal/repository"
	"github.com/leandro-lugaresi/hub"
	"github.com/mitchellh/mapstructure"

	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

type eventHandler func(w *Service, ev hub.Message)

type wireguardHelper struct {
	repo repository.Repository
}

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

	wgHelper := wireguardHelper{
		repo: w.repo,
	}

	peer, err := wgHelper.parseMessageForPeer(msg.DeviceID)

	if err != nil {
		w.logger.Warn(fmt.Sprintf("error while parsing deviceCreatedMessage - %s", err))
		return
	}

	if err := w.addPeer(peer); err != nil {
		w.logger.Warn(fmt.Sprintf("error while adding peer %s", err))
	}
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

	if err = w.deletePeer(&peer); err != nil {
		w.logger.Warn(fmt.Sprintf("error while deleting peer %s", err))
	}
}

func deviceUpdatedHandler(w *Service, ev hub.Message) {
	var msg event.DeviceUpdatedEvent
	err := mapstructure.Decode(ev.Fields, &msg)
	if err != nil {
		w.logger.Warn("error decode interface to struct")
		return
	}

	wgHelper := wireguardHelper{
		repo: w.repo,
	}

	peer, err := wgHelper.parseMessageForPeer(msg.DeviceID)

	if err != nil {
		w.logger.Warn(fmt.Sprintf("error while parsing deviceUpdatedMessage - %s", err))
		return
	}

	if err = w.updatePeer(peer); err != nil {
		w.logger.Warn(fmt.Sprintf("error while updating wg peer %s", err))
	}
}

func (w *wireguardHelper) parseMessageForPeer(deviceID uint) (*Peer, error) {
	device, err := w.repo.GetDevice(deviceID)

	if err != nil {
		return nil, fmt.Errorf("extracting of device infos failed %s", err)
	}

	publicKey, err := wgtypes.ParseKey(device.PublicKey)

	if err != nil {
		return nil, fmt.Errorf("error parsing wireguard public key %s", err)
	}

	peer := Peer{
		PublicKey:   publicKey,
		IPV4Address: net.ParseIP(device.IPv4Address),
		IPV6Address: net.ParseIP(device.IPv6Address),
	}

	return &peer, nil
}
