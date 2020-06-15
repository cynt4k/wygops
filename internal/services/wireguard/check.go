package wireguard

import (
	"fmt"
	"math/rand"
	"net"
	"regexp"

	"github.com/cynt4k/wygops/internal/models"
	"github.com/dspinhirne/netaddr-go"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

func (w *Service) checkV4Subnet(ip string) error {
	if ip == "" {
		return nil
	}
	server, err := netaddr.ParseIPv4Net(w.subnet.V4)
	if err != nil {
		return err
	}

	client, err := netaddr.ParseIPv4(ip)

	if err != nil {
		return err
	}

	if !server.Contains(client) {
		return ErrWrongSubnet
	}
	return nil
}

func (w *Service) checkV6Subnet(ip string) error {
	if ip == "" {
		return nil
	}
	server, err := netaddr.ParseIPv6Net(w.subnet.V6)
	if err != nil {
		return err
	}
	client, err := netaddr.ParseIPv6(ip)
	if err != nil {
		return err
	}

	if !server.Contains(client) {
		return ErrWrongSubnet
	}
	return nil
}

func (w *Service) getIPV4Subnet() ([]string, error) {
	inc := func(ipInc net.IP) {
		for j := len(ipInc) - 1; j >= 0; j-- {
			ipInc[j]++
			if ipInc[j] > 0 {
				break
			}
		}
	}
	ip, ipnet, err := net.ParseCIDR(w.subnet.V4)

	if err != nil {
		return nil, err
	}

	var ips []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}

	lenIPs := len(ips)
	switch {
	case lenIPs < 2:
		return ips, nil
	default:
		return ips[1 : len(ips)-1], nil
	}
}

func (w *Service) getAvailableIPV4(devices []*models.Device) (*net.IP, error) {
	allIps, err := w.getIPV4Subnet()

	if err != nil {
		return nil, err
	}

	if len(allIps) == 0 {
		return nil, ErrNoIPAvailable
	}

	var selectedIP *net.IP
	if len(devices) == 0 {
		for _, availableIP := range allIps {
			netIP := net.ParseIP(availableIP)
			if netIP.Equal(net.ParseIP(w.subnet.GatewayV4)) {
				continue
			}
			selectedIP = &netIP
			break
		}

		if selectedIP == nil {
			return nil, ErrNoIPAvailable
		}
		return selectedIP, nil
	}

	for _, device := range devices {
		ip := net.ParseIP(device.IPv4Address)
		for i, availableIP := range allIps {
			netIP := net.ParseIP(availableIP)
			if ip.Equal(netIP) {
				continue
			}
			if netIP.Equal(net.ParseIP(w.subnet.GatewayV4)) {
				continue
			}

			selectedIP = &netIP
			allIps = allIps[i:]
			break
		}
	}

	if selectedIP == nil {
		return nil, ErrNoIPAvailable
	}
	return selectedIP, nil
}

func (w *Service) getAvailableIPV6(devices []*models.Device) (*net.IP, error) {

	genAddress := func(sn string, retries int) (net.IP, error) {
		re, _ := regexp.Compile("(^.*)\\/.*")

		matches := re.FindStringSubmatch(sn)

		if len(matches) != 2 {
			return nil, ErrNoIPAvailable
		}

		subnet, err := netaddr.ParseIPv6(matches[1])

		if err != nil {
			return nil, err
		}

		var address *netaddr.IPv6
		for i := 0; i < retries; i++ {
			address = netaddr.NewIPv6(subnet.NetId(), rand.Uint64())

			if address.String() == w.subnet.GatewayV6 {
				err = ErrNoIPAvailable
				continue
			}
			break
		}

		if err != nil {
			return nil, err
		}

		w.logger.Info(address.String())
		return net.ParseIP(address.String()), nil
	}

	var ipv6 net.IP
	var err error

	ipv6, err = genAddress(w.subnet.V6, 10)

	if err != nil {
		return nil, err
	}

	if ipv6 == nil {
		return nil, fmt.Errorf("could not generate an ipv6 address - possible to small subnet")
	}

	if len(devices) == 0 {
		return &ipv6, nil
	}

	for _, device := range devices {
		var found bool
		for i := 0; i < 2; i++ {
			deviceAddress := net.ParseIP(device.IPv6Address)
			if deviceAddress.Equal(ipv6) {
				ipv6, err = genAddress(w.subnet.V6, 10)
				if err != nil {
					return nil, err
				}
				continue
			}
			found = true
			break
		}
		if found {
			return &ipv6, nil
		}
	}
	return nil, ErrNoIPAvailable
}

func (w *Service) parsePeer(device *Peer) (peer *wgtypes.PeerConfig, err error) {
	if err := w.checkV4Subnet(device.IPV4Address.String()); err != nil {
		return nil, err
	}
	if err := w.checkV6Subnet(device.IPV6Address.String()); err != nil {
		return nil, err
	}

	_, networkV4, err := net.ParseCIDR(fmt.Sprintf("%s/32", device.IPV4Address.String()))
	if err != nil {
		return nil, err
	}
	_, networkV6, err := net.ParseCIDR(fmt.Sprintf("%s/128", device.IPV6Address.String()))
	if err != nil {
		return nil, err
	}
	peer = &wgtypes.PeerConfig{
		PublicKey:  device.PublicKey,
		AllowedIPs: []net.IPNet{*networkV4, *networkV6},
	}
	return
}
