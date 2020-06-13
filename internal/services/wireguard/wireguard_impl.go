package wireguard

import (
	"fmt"
	"net"

	"github.com/cynt4k/wygops/cmd/config"
	"github.com/cynt4k/wygops/internal/repository"
	"github.com/leandro-lugaresi/hub"
	"go.uber.org/zap"
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

// Service : Service struct for wireguard
type Service struct {
	client *wgctrl.Client
	server *wgtypes.Device
	repo   repository.Repository
	hub    *hub.Hub
	config *config.Wireguard
	subnet *config.GeneralSubnet
	logger *zap.Logger
	peers  []*wgtypes.Peer
}

var (
	service *Service
)

// NewService : Create a new wireguard service
func NewService(repo repository.Repository, hub *hub.Hub, logger *zap.Logger, config *config.Config) (Wireguard, error) {
	if service != nil {
		return service, nil
	}
	wgConfig := config.Wireguard
	client, err := wgctrl.New()

	if err != nil {
		return nil, err
	}

	device, err := client.Device(wgConfig.Interface)

	if err != nil {
		return nil, err
	}

	wg := &Service{
		client: client,
		server: device,
		repo:   repo,
		hub:    hub,
		logger: logger,
		config: &config.Wireguard,
		subnet: &config.General.Subnet,
	}
	service = wg

	err = service.init()
	service.initEventhandler()

	if err != nil {
		return nil, err
	}

	return service, nil
}

func (w *Service) initEventhandler() {
	go func() {
		topics := make([]string, 0, len(handlerMap))
		for k := range handlerMap {
			topics = append(topics, k)
		}
		for msg := range w.hub.Subscribe(200, topics...).Receiver {
			h, ok := handlerMap[msg.Topic()]
			if ok {
				go h(w, msg)
			}
		}
	}()
}

func (w *Service) init() error {
	devices, err := w.repo.GetDevices()

	if err != nil {
		return err
	}

	for _, device := range devices {

		publicKey, err := wgtypes.ParseKey(device.PublicKey)

		if err != nil {
			w.logger.Warn(fmt.Sprintf("public key for device %v not valid - disable peer", device.ID))
			continue
		}

		wgDevice := &Peer{
			PublicKey:   publicKey,
			IPV4Address: net.ParseIP(device.IPv4Address),
			IPV6Address: net.ParseIP(device.IPv6Address),
		}

		if err := w.addPeer(wgDevice); err != nil {
			w.logger.Warn(err.Error())
		}
	}

	return nil
}

func (w *Service) addPeer(device *Peer) error {
	peer, err := w.parsePeer(device)

	if err != nil {
		return err
	}

	return w.client.ConfigureDevice(w.config.Interface, wgtypes.Config{
		Peers: []wgtypes.PeerConfig{*peer},
	})
}

func (w *Service) updatePeer(device *Peer) error {
	_, err := w.parsePeer(device)

	if err != nil {
		return err
	}

	var selectedPeer *wgtypes.Peer
	for _, peer := range w.server.Peers {
		if peer.PublicKey.String() == device.PublicKey.String() {
			selectedPeer = &peer
		}
	}

	if selectedPeer == nil {
		return w.addPeer(device)
	}

	var allPeers []wgtypes.PeerConfig
	for _, peer := range w.server.Peers {
		var config wgtypes.PeerConfig
		if peer.PublicKey.String() == selectedPeer.PublicKey.String() {
			config = wgtypes.PeerConfig{
				PublicKey:  selectedPeer.PublicKey,
				AllowedIPs: selectedPeer.AllowedIPs,
			}
		} else {
			config = wgtypes.PeerConfig{
				PublicKey:  peer.PublicKey,
				AllowedIPs: peer.AllowedIPs,
			}
		}
		allPeers = append(allPeers, config)
	}

	return w.client.ConfigureDevice(w.config.Interface, wgtypes.Config{
		Peers:        allPeers,
		ReplacePeers: true,
	})
}

func (w *Service) deletePeer(device *Peer) error {
	peer, err := w.parsePeer(device)

	if err != nil {
		return err
	}

	return w.client.ConfigureDevice(w.config.Interface, wgtypes.Config{
		Peers: []wgtypes.PeerConfig{
			wgtypes.PeerConfig{
				PublicKey: peer.PublicKey,
				Remove:    true,
			},
		},
	})

}

// CreatePeer : Create a new wireguard device
func (w *Service) CreatePeer() (*Peer, error) {
	privateKey, err := wgtypes.GeneratePrivateKey()

	if err != nil {
		return nil, err
	}
	publicKey := privateKey.PublicKey()
	devices, err := w.repo.GetDevices()

	if err != nil {
		return nil, err
	}

	ipv4, err := w.getAvailableIPV4(devices)
	if err != nil {
		return nil, err
	}

	ipv6, err := w.getAvailableIPV6(devices)

	if err != nil {
		return nil, err
	}

	peer := Peer{
		PrivateKey:  privateKey,
		PublicKey:   publicKey,
		IPV4Address: *ipv4,
		IPV6Address: *ipv6,
	}

	return &peer, nil
}
