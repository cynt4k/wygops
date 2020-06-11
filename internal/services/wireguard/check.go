package wireguard

import (
	"fmt"

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

	client, err := netaddr.ParseIPv4Net(ip)

	if err != nil {
		return err
	}

	result, err := server.Cmp(client)

	if err != nil {
		return err
	}

	if result != 0 {
		return fmt.Errorf("ipv4 of client not the same subnet")
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
	client, err := netaddr.ParseIPv6Net(ip)
	if err != nil {
		return err
	}

	result, err := server.Cmp(client)

	if err != nil {
		return err
	}

	if result != 0 {
		return fmt.Errorf("ipv6 of client not in the same subnet")
	}
	return nil
}
