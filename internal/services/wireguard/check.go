package wireguard

import (
	"fmt"
	"net"

	"github.com/cynt4k/wygops/internal/models"
	"github.com/dspinhirne/netaddr-go"
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

func (w *Service) getIPV4Subnet() ([]*net.IP, error) {
	inc := func(ip net.IP) {
		for j := len(ip) - 1; j >= 0; j-- {
			ip[j]++
			if ip[j] > 0 {
				break
			}
		}
	}
	ip, ipnet, err := net.ParseCIDR(w.subnet.V4)

	if err != nil {
		return nil, err
	}

	var ips []*net.IP
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, &ip)
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
	for _, device := range devices {
		ip := net.ParseIP(device.IPv4Address)
		for _, availableIP := range allIps {
			if ip.Equal(*availableIP) {
				continue
			}
			return availableIP, nil
		}
	}
	return nil, ErrNoIPAvailable
}

func (w *Service) getAvailableIPV6(devices []*models.Device) (*net.IP, error) {

	genAddress := func(sn string, retries int) (net.IP, error) {
		subnet, err := netaddr.ParseIPv6(sn)

		if err != nil {
			return nil, err
		}

		mask, err := netaddr.NewMask128(128)

		if err != nil {
			return nil, err
		}

		var address netaddr.IPv6
		for i := 0; i < retries; i++ {
			address, err := netaddr.NewIPv6Net(subnet, mask)

			if err != nil {
				return nil, err
			}

			if address.String() == w.subnet.GatewayV6 {
				err = ErrNoIPAvailable
				continue
			}
		}

		if err != nil {
			return nil, err
		}

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
