package wireguard

import (
	"fmt"
	"net"

	"github.com/cynt4k/wygops/cmd/config"
	"github.com/cynt4k/wygops/internal/repository"
	"go.uber.org/zap"
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

// Service : Service struct for wireguard
type Service struct {
	client *wgctrl.Client
	server *wgtypes.Device
	repo   repository.Repository
	config *config.Wireguard
	subnet *config.GeneralSubnet
	logger *zap.Logger
	peers  []*wgtypes.Peer
}

var (
	service *Service
)

// NewService : Create a new wireguard service
func NewService(repo repository.Repository, logger *zap.Logger, config *config.Config) (Wireguard, error) {
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
		logger: logger,
		config: &config.Wireguard,
		subnet: &config.General.Subnet,
	}
	service = wg

	err = service.init()

	if err != nil {
		return nil, err
	}

	return service, nil
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

	if err := w.checkV4Subnet(device.IPV4Address.String()); err != nil {
		return err
	}
	if err := w.checkV6Subnet(device.IPV6Address.String()); err != nil {
		return err
	}

	_, networkV4, err := net.ParseCIDR(fmt.Sprintf("%s/32", device.IPV4Address.String()))
	if err != nil {
		return err
	}
	_, networkV6, err := net.ParseCIDR(fmt.Sprintf("%s/128", device.IPV6Address.String()))
	if err != nil {
		return err
	}

	peer := &wgtypes.PeerConfig{
		PublicKey:  device.PublicKey,
		AllowedIPs: []net.IPNet{*networkV4, *networkV6},
	}

	err = w.client.ConfigureDevice(w.config.Interface, wgtypes.Config{
		Peers: []wgtypes.PeerConfig{*peer},
	})

	if err != nil {
		return err
	}
	return nil
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

	peer := Peer{
		PrivateKey:  privateKey,
		PublicKey:   publicKey,
		IPV4Address: *ipv4,
		IPV6Address: *ipv6,
	}

	return &peer, nil
}
