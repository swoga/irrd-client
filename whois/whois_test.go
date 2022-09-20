package whois

import (
	"fmt"
	"testing"
	"time"

	"github.com/swoga/irrd-client/mock"
)

func TestEnableMultiCommand(t *testing.T) {
	mc, s := mock.SetupMockConn()
	c := NewFromBufferedConn(mc)

	go func() {
		defer s.Close()
		s.Read(t, "!!\n")
	}()

	err := c.EnableMultiCommand()
	if err != nil {
		t.Fatal(err)
	}
}

func TestSetIdleTimeout(t *testing.T) {

	mc, s := mock.SetupMockConn()
	c := NewFromBufferedConn(mc)

	secs := 300

	go func() {
		defer s.Close()
		s.Read(t, "!t"+fmt.Sprint(secs)+"\n")
		s.Write(t, "C\n")
	}()

	err := c.SetIdleTimout(time.Second * 300)
	if err != nil {
		t.Fatal(err)
	}
}
