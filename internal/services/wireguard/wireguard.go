package wireguard

import (
	"net"

	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

// Peer : Wireguard Peer
type Peer struct {
	PrivateKey  wgtypes.Key
	PublicKey   wgtypes.Key
	IPV4Address net.IP
	IPV6Address net.IP
}

// Wireguard : Wireguard service interface
type Wireguard interface {
	CreatePeer() (*Peer, error)
}
