package whois

import (
	"fmt"
	"strings"

	"inet.af/netaddr"
)

func (w whois) GetRoutesByOrigin(p IPProto, asn uint32) ([]netaddr.IPPrefix, error) {
	query := ""
	switch p {
	case IP4:
		query = "g"
	case IP6:
		query = "6"
	case IPany:
		return nil, fmt.Errorf("query does not support any proto")
	default:
		return nil, fmt.Errorf("unknown proto: %v", p)
	}
	str, err := w.Query("!" + query + "AS" + fmt.Sprintf("%v", asn))
	if err != nil {
		return nil, err
	}
	if str == "" {
		return make([]netaddr.IPPrefix, 0), nil
	}
	strNets := strings.Split(str, " ")
	return parseNetworks(strNets)
}
