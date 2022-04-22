package whois

import (
	"fmt"
	"testing"
	"time"

	"github.com/swoga/irrd-client/asynctest"
	"github.com/swoga/irrd-client/mock"
)

func TestEnableMultiCommand(t *testing.T) {
	mc, s := mock.SetupMockConn()
	c := NewFromBufferedConn(mc)

	asynctest.New(t,
		func(t asynctest.T) {
			s.Read(t, "!!\n")
		},
		func(t asynctest.T) {
			err := c.EnableMultiCommand()
			if err != nil {
				t.Fatal(err)
			}
		})
}

func TestSetIdleTimeout(t *testing.T) {

	mc, s := mock.SetupMockConn()
	c := NewFromBufferedConn(mc)

	secs := 300

	asynctest.New(t,
		func(t asynctest.T) {
			s.Read(t, "!t"+fmt.Sprint(secs)+"\n")
			s.Write(t, "C\n")
		},
		func(t asynctest.T) {
			err := c.SetIdleTimout(time.Second * 300)
			if err != nil {
				t.Fatal(err)
			}
		})
}
