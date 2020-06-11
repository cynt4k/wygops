package wireguard

import (
	"github.com/cynt4k/wygops/cmd/config"
	"github.com/cynt4k/wygops/internal/repository"
	"golang.zx2c4.com/wireguard/wgctrl"
)

// Service : Service struct for wireguard
type Service struct {
	client     *wgctrl.Client
	repo       repository.Repository
	config     *config.Wireguard
	subnet     *config.GeneralSubnet
	interfaces []*Device
}

// NewService : Create a new wireguard service
func NewService(repo repository.Repository, wgConfig *config.Wireguard, subnet *config.GeneralSubnet) (Wireguard, error) {
	client, err := wgctrl.New()

	if err != nil {
		return nil, err
	}
	wg := &Service{
		client: client,
		repo:   repo,
		config: wgConfig,
		subnet: subnet,
	}

	return wg, nil
}

func (w *Service) init() error {
	devices, err := w.repo.GetDevices()

	if err != nil {
		return err
	}

	for _, device := range devices {

		if err = w.checkV4Subnet(device.IPv4Address); err != nil {
			return err
		}

		if err = w.checkV6Subnet(device.IPv6Address); err != nil {
			return err
		}
	}

	return nil
}

// CreateDevice : Create a new wireguard device
func (w *Service) CreateDevice(name string) (*Device, error) {
	// device := w.client.ConfigureDevice(name)
	return nil, nil
}

// DeleteDevice : Delete a wireguard device
func (w *Service) DeleteDevice(device *Device) error {
	return nil
}

// UpdateDevice : Update a wireguard device
func (w *Service) UpdateDevice(device *Device) (*Device, error) {
	return nil, nil
}
