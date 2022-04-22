package async

import "github.com/swoga/irrd-client/whois"

func (a async) GetAsSetMembersRecrusive(set string) AsSetMembers {
	rch := make(chan Result[[]uint32], 1)
	a.queries <- func(w whois.Whois) error {
		s, err := w.GetAsSetMembersRecrusive(set)
		rch <- result[[]uint32]{s, err}
		return err
	}
	return rch
}

func (a async) GetSetMembers(set string, recursive bool) SetMembers {
	rch := make(chan Result[[]string], 1)
	a.queries <- func(w whois.Whois) error {
		s, err := w.GetSetMembers(set, recursive)
		rch <- result[[]string]{s, err}
		return err
	}
	return rch
}
