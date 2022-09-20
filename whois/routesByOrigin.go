package whois

import (
	"fmt"
	"net/netip"
	"strings"
)

func (w whois) GetRoutesByOrigin(p IPProto, asn uint32) ([]netip.Prefix, error) {
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
		return make([]netip.Prefix, 0), nil
	}
	strNets := strings.Split(str, " ")
	return parseNetworks(strNets)
}
