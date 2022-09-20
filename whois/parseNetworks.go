package whois

import (
	"fmt"
	"net/netip"
)

func parseNetworks(strNets []string) ([]netip.Prefix, error) {
	nets := make([]netip.Prefix, 0, len(strNets))
	for _, strNet := range strNets {
		net, err := netip.ParsePrefix(strNet)
		if err != nil {
			return nil, fmt.Errorf("cannot parse %v into network: %v", strNet, err)
		}
		nets = append(nets, net)
	}
	return nets, nil
}
