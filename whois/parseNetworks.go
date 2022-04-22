package whois

import (
	"fmt"

	"inet.af/netaddr"
)

func parseNetworks(strNets []string) ([]netaddr.IPPrefix, error) {
	nets := make([]netaddr.IPPrefix, 0, len(strNets))
	for _, strNet := range strNets {
		net, err := netaddr.ParseIPPrefix(strNet)
		if err != nil {
			return nil, fmt.Errorf("cannot parse %v into network: %v", strNet, err)
		}
		nets = append(nets, net)
	}
	return nets, nil
}
