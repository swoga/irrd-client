package whois

import (
	"fmt"
	"net/netip"
	"strings"
)

func (w whois) GetRoutesBySet(p IPProto, set string) ([]netip.Prefix, error) {
	query := ""
	switch p {
	case IP4:
		query = "4"
	case IP6:
		query = "6"
	case IPany:
		query = ""
	default:
		return nil, fmt.Errorf("unknown proto: %v", p)
	}

	str, err := w.Query("!a" + query + set)
	if err != nil {
		return nil, err
	}
	if str == "" {
		return make([]netip.Prefix, 0), nil
	}
	strNets := strings.Split(str, " ")
	return parseNetworks(strNets)
}
