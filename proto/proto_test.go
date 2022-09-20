package proto_test

import (
	"fmt"
	"testing"

	"github.com/swoga/irrd-client/mock"
)

func TestClientA(t *testing.T) {
	c, s := mock.SetupMockClient()

	query := "ababab"
	want := "cdcdcd"

	go func() {
		defer s.Close()
		s.Read(t, query+"\n")
		s.Write(t, "A"+fmt.Sprintf("%d", len(want)+1)+"\n"+want+"\nC\n")
	}()

	has, err := c.Query(query)
	if err != nil {
		t.Fatal(err)
	}
	if has != want {
		t.Fatalf("client got: %s, want: %s", has, want)
	}
}

func TestClientF(t *testing.T) {
	c, s := mock.SetupMockClient()

	query := "ababab"
	want := "<error message>"

	go func() {
		defer s.Close()
		s.Read(t, query+"\n")
		s.Write(t, "F "+want+"\n")
	}()

	_, err := c.Query(query)
	if err == nil {
		t.Fatal("expected error")
	}

	if err.Error() != "server returned error: "+want {
		t.Fatalf("client got error: %s, want: %s", err, want)
	}
}

func TestClientX(t *testing.T) {
	c, s := mock.SetupMockClient()

	query := "ababab"
	want := "c"

	go func() {
		defer s.Close()
		s.Read(t, query+"\n")
		s.Write(t, want+"\n")
	}()

	_, err := c.Query(query)
	if err == nil {
		t.Fatal("expected error")
	}

	if err.Error() != "unknown response: "+want {
		t.Fatalf("client got error: %s, want: %s", err, want)
	}

}
