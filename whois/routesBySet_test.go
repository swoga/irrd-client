package whois

import (
	"fmt"
	"net/netip"
	"testing"
	"time"

	"github.com/swoga/irrd-client/mock"
)

func TestGetRoutesBySet(t *testing.T) {
	mc, s := mock.SetupMockConn()
	c := NewFromBufferedConn(mc)

	set := "AS-TEST"
	want := []netip.Prefix{netip.MustParsePrefix("10.0.0.0/8"), netip.MustParsePrefix("192.168.0.0/16")}
	raw := "10.0.0.0/8 192.168.0.0/16"

	go func() {
		defer s.Close()
		s.Read(t, "!a4"+set+"\n")
		s.Write(t, "A"+fmt.Sprintf("%d", len(raw)+1)+"\n"+raw+"\nC\n")
	}()

	has, err := c.GetRoutesBySet(IP4, set)
	if err != nil {
		t.Fatal(err)
	}

	if len(has) != len(want) {
		t.Fatalf("client got n results: %v, want: %v", len(has), len(want))
	}

	for i := range has {
		if has[i] != want[i] {
			t.Fatalf("client got: %v, want: %v", has[i], want[i])
		}
	}
}

func TestGetRoutesBySetEmpty(t *testing.T) {
	mc, s := mock.SetupMockConn()
	c := NewFromBufferedConn(mc)

	set := "AS-TEST"

	go func() {
		defer s.Close()
		s.Read(t, "!a4"+set+"\n")
		s.Write(t, "D\nC\n")
	}()

	has, err := c.GetRoutesBySet(IP4, set)
	if err != nil {
		t.Fatal(err)
	}

	if len(has) != 0 {
		t.Fatalf("client got n results: %v, want: %v", len(has), 0)
	}
}

func TestLiveGetRoutes4BySet(t *testing.T) {
	c, err := New("rr.ntt.net:43", time.Second*10)
	if err != nil {
		t.Fatal(err)
	}

	x, err := c.GetRoutesBySet(IPany, "AS-TEST")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(x)
}
