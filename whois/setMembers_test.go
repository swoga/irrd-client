package whois

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/swoga/irrd-client/mock"
)

func TestGetSetMembers(t *testing.T) {
	mc, s := mock.SetupMockConn()
	c := NewFromBufferedConn(mc)

	query := "AS-TEST"
	want := []uint32{1, 2}
	raw := "AS1 AS2"

	go func() {
		defer s.Close()
		s.Read(t, "!i"+query+",1\n")
		s.Write(t, "A"+fmt.Sprintf("%d", len(raw)+1)+"\n"+raw+"\nC\n")
	}()

	has, err := c.GetAsSetMembersRecrusive(query)
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

func TestGetSetMembersEmpty(t *testing.T) {
	mc, s := mock.SetupMockConn()
	c := NewFromBufferedConn(mc)

	query := "AS-TEST"

	go func() {
		defer s.Close()
		s.Read(t, "!i"+query+",1\n")
		s.Write(t, "D\nC\n")
	}()

	has, err := c.GetAsSetMembersRecrusive(query)
	if err != nil {
		t.Fatal(err)
	}

	if len(has) != 0 {
		t.Fatalf("client got n results: %v, want: %v", len(has), 0)
	}
}

func TestGetAsSetMembersInvalidInt(t *testing.T) {
	mc, s := mock.SetupMockConn()
	c := NewFromBufferedConn(mc)

	query := "AS-TEST"

	go func() {
		defer s.Close()
		s.Read(t, "!i"+query+",1\n")
		res := "AS100 ASXXX"
		s.Write(t, "A"+fmt.Sprint(len(res)+1)+"\n"+res+"\nC\n")
	}()

	_, err := c.GetAsSetMembersRecrusive(query)
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.HasPrefix(err.Error(), "cannot parse ASXXX into ASN:") {
		t.Fatalf("expected other error, got: %v", err)
	}
}

func TestGetAsSetMembersInvalidShort(t *testing.T) {
	mc, s := mock.SetupMockConn()
	c := NewFromBufferedConn(mc)

	query := "AS-TEST"

	go func() {
		defer s.Close()
		s.Read(t, "!i"+query+",1\n")
		res := "AS100 X"
		s.Write(t, "A"+fmt.Sprint(len(res)+1)+"\n"+res+"\nC\n")
	}()

	_, err := c.GetAsSetMembersRecrusive(query)
	if err == nil {
		t.Fatalf("expected error")
	}
	if err.Error() != "cannot parse X into ASN: too short" {
		t.Fatalf("expected other error, got: %v", err)
	}
}

func TestLiveGetAsSetMembersRecrusive(t *testing.T) {
	c, err := New("whois.radb.net:43", time.Second*10)
	if err != nil {
		t.Fatal(err)
	}
	x, err := c.GetAsSetMembersRecrusive("AS-TEST")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(x)
}
