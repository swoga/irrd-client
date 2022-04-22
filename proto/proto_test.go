package proto_test

import (
	"fmt"
	"testing"

	"github.com/swoga/irrd-client/asynctest"
	"github.com/swoga/irrd-client/mock"
)

func TestClientA(t *testing.T) {
	c, s := mock.SetupMockClient()

	query := "ababab"
	want := "cdcdcd"

	asynctest.New(t,
		func(t asynctest.T) {
			s.Read(t, query+"\n")
			s.Write(t, "A"+fmt.Sprintf("%d", len(want)+1)+"\n"+want+"\nC\n")
		},
		func(t asynctest.T) {
			has, err := c.Query(query)
			if err != nil {
				t.Fatal(err)
			}
			if has != want {
				t.Fatalf("client got: %s, want: %s", has, want)
			}
		})
}

func TestClientF(t *testing.T) {
	c, s := mock.SetupMockClient()

	query := "ababab"
	want := "<error message>"

	asynctest.New(t,
		func(t asynctest.T) {
			s.Read(t, query+"\n")
			s.Write(t, "F "+want+"\n")
		},
		func(t asynctest.T) {
			_, err := c.Query(query)
			if err == nil {
				t.Fatal("expected error")
			}

			if err.Error() != "server returned error: "+want {
				t.Fatalf("client got error: %s, want: %s", err, want)
			}
		})
}

func TestClientX(t *testing.T) {
	c, s := mock.SetupMockClient()

	query := "ababab"
	want := "c"

	asynctest.New(t,
		func(t asynctest.T) {
			s.Read(t, query+"\n")
			s.Write(t, want+"\n")
		},
		func(t asynctest.T) {
			_, err := c.Query(query)
			if err == nil {
				t.Fatal("expected error")
			}

			if err.Error() != "unknown response: "+want {
				t.Fatalf("client got error: %s, want: %s", err, want)
			}
		})

}
