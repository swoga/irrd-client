package async

import (
	"github.com/swoga/irrd-client/whois"
	"inet.af/netaddr"
)

func (a async) GetRoutesByOrigin(p whois.IPProto, asn uint32) Routes {
	rch := make(chan Result[[]netaddr.IPPrefix], 1)
	a.queries <- func(w whois.Whois) error {
		defer close(rch)
		p, err := w.GetRoutesByOrigin(p, asn)
		rch <- result[[]netaddr.IPPrefix]{p, err}
		return err
	}
	return rch
}
