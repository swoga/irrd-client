package async

import (
	"strings"

	"github.com/swoga/irrd-client/whois"
)

func (a async) GetVersion() <-chan Result[string] {
	rch := make(chan Result[string], 1)
	a.queries <- func(w whois.Whois) error {
		defer close(rch)
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
