package async

import (
	"github.com/swoga/irrd-client/whois"
	"inet.af/netaddr"
)

func (a async) GetRoutesBySet(p whois.IPProto, set string) Routes {
	rch := make(chan Result[[]netaddr.IPPrefix], 1)
	a.queries <- func(w whois.Whois) error {
		p, err := w.GetRoutesBySet(p, set)
		rch <- result[[]netaddr.IPPrefix]{p, err}
		return err
	}
	return rch
}
