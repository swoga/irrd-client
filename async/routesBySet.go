package async

import (
	"net/netip"

	"github.com/swoga/irrd-client/whois"
)

func (a async) GetRoutesBySet(p whois.IPProto, set string) Routes {
	rch := make(chan Result[[]netip.Prefix], 1)
	a.queries <- func(w whois.Whois) error {
		defer close(rch)
		p, err := w.GetRoutesBySet(p, set)
		rch <- result[[]netip.Prefix]{p, err}
		return err
	}
	return rch
}
