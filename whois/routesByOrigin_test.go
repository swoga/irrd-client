package whois

import (
	"fmt"
	"testing"
	"time"

	"github.com/swoga/irrd-client/mock"
	"inet.af/netaddr"
)

func TestGetRoutesByOrigin(t *testing.T) {
	mc, s := mock.SetupMockConn()
	c := NewFromBufferedConn(mc)

	var asn uint32 = 1234
	asnStr := fmt.Sprintf("%d", asn)
	want := []netaddr.IPPrefix{netaddr.MustParseIPPrefix("10.0.0.0/8"), netaddr.MustParseIPPrefix("192.168.0.0/16")}
	raw := "10.0.0.0/8 192.168.0.0/16"

	go func() {
		defer s.Close()
		s.Read(t, "!gAS"+asnStr+"\n")
		s.Write(t, "A"+fmt.Sprintf("%d", len(raw)+1)+"\n"+raw+"\nC\n")
	}()

	has, err := c.GetRoutesByOrigin(IP4, asn)
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

func TestGetRoutesByOriginEmpty(t *testing.T) {
	mc, s := mock.SetupMockConn()
	c := NewFromBufferedConn(mc)

	var asn uint32 = 1234
	asnStr := fmt.Sprintf("%d", asn)

	go func() {
		defer s.Close()
		s.Read(t, "!gAS"+asnStr+"\n")
		s.Write(t, "D\nC\n")
	}()

	has, err := c.GetRoutesByOrigin(IP4, asn)
	if err != nil {
		t.Fatal(err)
	}

	if len(has) != 0 {
		t.Fatalf("client got n results: %v, want: %v", len(has), 0)
	}
}

func TestLiveGetRoutes4ByOrigin(t *testing.T) {
	c, err := New("whois.radb.net:43", time.Second*10)
	if err != nil {
		t.Fatal(err)
	}

	x, err := c.GetRoutesByOrigin(IP4, 1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(x)
}
