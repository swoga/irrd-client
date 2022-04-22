package async

import (
	"time"

	"github.com/swoga/irrd-client/whois"
	"inet.af/netaddr"
)

type Whois interface {
	Start(n int) error
	Stop()
	// send !t<timeout>
	SetIdleTimout(t time.Duration)
	// send query !i<set-name> or !i<set-name>,1 depending on bool, safe for concurrent use
	GetSetMembers(set string, recursive bool) SetMembers
	// send query !i<set-name>,1, safe for concurrent use
	GetAsSetMembersRecrusive(set string) AsSetMembers
	// send query !gAS<asn> or !6AS<asn> depending on p, safe for concurrent use
	GetRoutesByOrigin(p whois.IPProto, asn uint32) Routes
	// send query !a<as-set-name> !a4<as-set-name> or !a6<as-set-name> depending on p, safe for concurrent use
	GetRoutesBySet(p whois.IPProto, set string) Routes
}

type SetMembers = <-chan Result[[]string]
type AsSetMembers = <-chan Result[[]uint32]
type Routes = <-chan Result[[]netaddr.IPPrefix]

type async struct {
	cache       whois.WhoisCache
	address     string
	timeout     time.Duration
	idleTimeout *time.Duration
	queries     chan func(whois.Whois) error
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
