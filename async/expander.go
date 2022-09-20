package async

import (
	"net/netip"
	"strconv"

	"github.com/swoga/irrd-client/whois"
)

func (a async) expandSetRoutes(p whois.IPProto, set string) Routes {
	rch := make(chan Result[[]netip.Prefix])

	go func() {
		defer close(rch)

		ci := a.GetSetMembers(set, true)
		r := <-ci

		if r.Error() != nil {
			rch <- result[[]netip.Prefix]{nil, r.Error()}
			return
		}

		var nets []netip.Prefix
		var originResults []<-chan Result[[]netip.Prefix]

		for _, member := range r.Data() {
			if len(member) > 2 && member[:2] == "AS" {
				asn, err := strconv.Atoi(member[2:])
				if err != nil {
					rch <- result[[]netip.Prefix]{nil, err}
					return
				}
				o := a.GetRoutesByOrigin(p, uint32(asn))
				originResults = append(originResults, o)
			} else {
				net, err := netip.ParsePrefix(member)
				if err != nil {
					rch <- result[[]netip.Prefix]{nil, err}
					return
				}
				is6 := p == whois.IP6
				if net.Addr().Is6() != is6 {
					continue
				}

				nets = append(nets, net)
			}
		}

		for _, r := range originResults {
			or := <-r
			if or.Error() != nil {
				rch <- result[[]netip.Prefix]{nil, or.Error()}
				return
			}
			nets = append(nets, or.Data()...)
		}

		rch <- result[[]netip.Prefix]{nets, nil}
	}()

	return rch
}

func (a async) ExpandSet(p whois.IPProto, set string) Routes {
	if !a.supportsBySet {
		return a.expandSetRoutes(p, set)
	}

	return a.GetRoutesBySet(p, set)
}
