package expander

import (
	"strconv"

	"github.com/swoga/irrd-client/async"
	"github.com/swoga/irrd-client/whois"
	"inet.af/netaddr"
)

func ExpandSetRoutes(w async.Whois, p whois.IPProto, set string) (*netaddr.IPSet, error) {
	ci := w.GetSetMembers(set, true)
	r := <-ci

	if r.Error() != nil {
		return nil, r.Error()
	}

	var nets netaddr.IPSetBuilder
	var originResults []<-chan async.Result[[]netaddr.IPPrefix]

	for _, member := range r.Data() {
		if len(member) > 2 && member[:2] == "AS" {
			asn, err := strconv.Atoi(member[2:])
			if err != nil {
				return nil, err
			}
			o := w.GetRoutesByOrigin(p, uint32(asn))
			originResults = append(originResults, o)
		} else {
			net, err := netaddr.ParseIPPrefix(member)
			if err != nil {
				return nil, err
			}
			is6 := p == whois.IP6
			if net.IP().Is6() != is6 {
				continue
			}

			nets.AddPrefix(net)
		}
	}

	for _, r := range originResults {
		or := <-r
		if or.Error() != nil {
			return nil, or.Error()
		}
		for _, p := range or.Data() {
			nets.AddPrefix(p)
		}
	}

	ipSet, err := nets.IPSet()
	if err != nil {
		return nil, err
	}

	return ipSet, nil
}
