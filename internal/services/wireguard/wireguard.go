package wireguard

// Device : Wireguard Device
type Device struct {
	PrivateKey  string
	PublicKey   string
	IPV4Address string
	IPV6Address string
}

// Wireguard : Wireguard service interface
type Wireguard interface {
	CreateDevice(name string) (*Device, error)
	DeleteDevice(device *Device) error
	UpdateDevice(device *Device) (*Device, error)
}
