package async

import (
	"strconv"
	"strings"
	"time"

	"github.com/swoga/irrd-client/whois"
	"inet.af/netaddr"
)

type Whois interface {
	Start(n int) error
	Stop()
	// send !t<timeout>
	SetIdleTimout(t time.Duration)
	// send !v
	GetVersion() <-chan Result[string]
	// send query !i<set-name> or !i<set-name>,1 depending on bool, safe for concurrent use
	GetSetMembers(set string, recursive bool) SetMembers
	// send query !i<set-name>,1, safe for concurrent use
	GetAsSetMembersRecrusive(set string) AsSetMembers
	// send query !gAS<asn> or !6AS<asn> depending on p, safe for concurrent use
	GetRoutesByOrigin(p whois.IPProto, asn uint32) Routes
	// send query !a<as-set-name> !a4<as-set-name> or !a6<as-set-name> depending on p, safe for concurrent use
	GetRoutesBySet(p whois.IPProto, set string) Routes

	ExpandSet(p whois.IPProto, set string) Routes
}

type SetMembers = <-chan Result[[]string]
type AsSetMembers = <-chan Result[[]uint32]
type Routes = <-chan Result[[]netaddr.IPPrefix]

type async struct {
	cache          whois.WhoisCache
	address        string
	timeout        time.Duration
	idleTimeout    *time.Duration
	queries        chan func(whois.Whois) error
	checkedVersion bool
	supportsBySet  bool
}

func New(address string, timeout time.Duration) Whois {
	return &async{
		cache:   whois.NewCache(),
		address: address,
		timeout: timeout,
		queries: make(chan func(whois.Whois) error),
	}
}

func (a *async) Start(n int) error {
	for i := 0; i < n; i++ {
		err := a.startWorker()
		if err != nil {
			return err
		}
	}

	err := a.checkVersion()
	if err != nil {
		return err
	}

	return nil
}

func (a async) Stop() {
	close(a.queries)
}

func (a *async) SetIdleTimout(t time.Duration) {
	a.idleTimeout = &t
}

func (a async) createWhois() (whois.Whois, error) {
	w, err := whois.New(a.address, a.timeout)
	if err != nil {
		return nil, err
	}
	err = w.EnableMultiCommand()
	if err != nil {
		return nil, err
	}
	if a.idleTimeout != nil {
		err = w.SetIdleTimout(*a.idleTimeout)
		if err != nil {
			return nil, err
		}
	}
	w = a.cache.UseCache(w)
	return w, nil
}

func (a async) startWorker() error {
	w, err := a.createWhois()
	if err != nil {
		return err
	}
	go a.worker(w)

	return nil
}

func (a async) worker(w whois.Whois) {
	for {
		q, ok := <-a.queries
		if !ok {
			// close the underlying connection if the channel got closed
			w.Close()
			break
		}
		err := q(w)
		// if any error occurs during the query, restart this worker
		if err != nil {
			a.startWorker()
			break
		}
	}
}

func (a async) GetVersion() <-chan Result[string] {
	rch := make(chan Result[string], 1)
	a.queries <- func(w whois.Whois) error {
		s, err := w.GetVersion()
		rch <- result[string]{s, err}
		return err
	}
	return rch
}

func (a *async) checkVersion() error {
	if a.checkedVersion {
		return nil
	}
	a.checkedVersion = true

	vch := a.GetVersion()
	r := <-vch
	if r.Error() != nil {
		return r.Error()
	}

	v := r.Data()
	header := "version "
	i := strings.Index(v, header)
	if i != -1 {
		start := i + len(header)
		version := v[start : start+1]
		a.supportsBySet = version == "4"
	}

	return nil
}

func (a async) expandSetRoutes(p whois.IPProto, set string) Routes {
	rch := make(chan Result[[]netaddr.IPPrefix])

	go func() {
		defer close(rch)

		ci := a.GetSetMembers(set, true)
		r := <-ci

		if r.Error() != nil {
			rch <- result[[]netaddr.IPPrefix]{nil, r.Error()}
			return
		}

		var nets []netaddr.IPPrefix
		var originResults []<-chan Result[[]netaddr.IPPrefix]

		for _, member := range r.Data() {
			if len(member) > 2 && member[:2] == "AS" {
				asn, err := strconv.Atoi(member[2:])
				if err != nil {
					rch <- result[[]netaddr.IPPrefix]{nil, err}
					return
				}
				o := a.GetRoutesByOrigin(p, uint32(asn))
				originResults = append(originResults, o)
			} else {
				net, err := netaddr.ParseIPPrefix(member)
				if err != nil {
					rch <- result[[]netaddr.IPPrefix]{nil, err}
					return
				}
				is6 := p == whois.IP6
				if net.IP().Is6() != is6 {
					continue
				}

				nets = append(nets, net)
			}
		}

		for _, r := range originResults {
			or := <-r
			if or.Error() != nil {
				rch <- result[[]netaddr.IPPrefix]{nil, or.Error()}
				return
			}
			nets = append(nets, or.Data()...)
		}

		rch <- result[[]netaddr.IPPrefix]{nets, nil}
	}()

	return rch
}

func (a async) ExpandSet(p whois.IPProto, set string) Routes {
	if !a.supportsBySet {
		return a.expandSetRoutes(p, set)
	}

	return a.GetRoutesBySet(p, set)
}
