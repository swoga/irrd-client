package async

import (
	"net/netip"

	"github.com/swoga/irrd-client/whois"
)

func (a async) GetRoutesByOrigin(p whois.IPProto, asn uint32) Routes {
	rch := make(chan Result[[]netip.Prefix], 1)
	a.queries <- func(w whois.Whois) error {
		defer close(rch)
		p, err := w.GetRoutesByOrigin(p, asn)
		rch <- result[[]netip.Prefix]{p, err}
		return err
	}
	return rch
}
