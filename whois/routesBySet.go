package whois

import (
	"fmt"
	"strings"

	"inet.af/netaddr"
)

func (w whois) GetRoutesBySet(p IPProto, set string) ([]netaddr.IPPrefix, error) {
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
		return make([]netaddr.IPPrefix, 0), nil
	}
	strNets := strings.Split(str, " ")
	return parseNetworks(strNets)
}
