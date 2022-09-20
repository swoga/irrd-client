package whois

import (
	"fmt"

	"net/netip"

	"github.com/swoga/irrd-client/cache"
)

type WhoisCache interface {
	UseCache(w Whois) Whois
}

func NewCache() WhoisCache {
	return whoisCache{
		setMembers:   cache.New[string, []string](),
		asSetMembers: cache.New[string, []uint32](),
		routes:       cache.New[string, []netip.Prefix](),
	}
}

type whoisCache struct {
	setMembers   cache.Cache[string, []string]
	asSetMembers cache.Cache[string, []uint32]
	routes       cache.Cache[string, []netip.Prefix]
}

func (wc whoisCache) UseCache(w Whois) Whois {
	return whoisCached{
		w,
		wc,
	}
}

type whoisCached struct {
	Whois
	cache whoisCache
}

func (wc whoisCached) GetRoutesByOrigin(p IPProto, asn uint32) ([]netip.Prefix, error) {
	key := fmt.Sprint(p) + fmt.Sprint(asn)

	value, found := wc.cache.routes.Get(key)
	if found {
		return value, nil
	}

	value, err := wc.Whois.GetRoutesByOrigin(p, asn)
	if err != nil {
		return nil, err
	}

	wc.cache.routes.Set(key, value)

	return value, nil
}

func (wc whoisCached) GetRoutesBySet(p IPProto, set string) ([]netip.Prefix, error) {
	key := fmt.Sprint(p) + set

	value, found := wc.cache.routes.Get(key)
	if found {
		return value, nil
	}

	value, err := wc.Whois.GetRoutesBySet(p, set)
	if err != nil {
		return nil, err
	}

	wc.cache.routes.Set(key, value)

	return value, nil
}

func (wc whoisCached) GetSetMembers(set string, recursive bool) ([]string, error) {
	flag := "-"
	if recursive {
		flag = "r"
	}
	key := flag + set

	value, found := wc.cache.setMembers.Get(key)
	if found {
		return value, nil
	}

	value, err := wc.Whois.GetSetMembers(set, recursive)
	if err != nil {
		return nil, err
	}

	wc.cache.setMembers.Set(key, value)

	return value, nil
}

func (wc whoisCached) GetAsSetMembersRecrusive(set string) ([]uint32, error) {
	key := set

	value, found := wc.cache.asSetMembers.Get(key)
	if found {
		return value, nil
	}

	value, err := wc.Whois.GetAsSetMembersRecrusive(set)
	if err != nil {
		return nil, err
	}

	wc.cache.asSetMembers.Set(key, value)

	return value, nil
}
