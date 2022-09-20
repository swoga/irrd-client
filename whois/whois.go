package whois

import (
	"fmt"
	"io"
	"net"
	"net/netip"
	"time"

	"github.com/swoga/irrd-client/conn"
	"github.com/swoga/irrd-client/proto"
)

type IPProto byte

const (
	IP4 IPProto = iota
	IP6
	IPany
)

type Whois interface {
	// send !!
	EnableMultiCommand() error
	// send !t<timeout>
	SetIdleTimout(d time.Duration) error
	// send !v
	GetVersion() (string, error)

	// send query !i<set-name> or !i<set-name>,1 depending on bool
	GetSetMembers(set string, recursive bool) ([]string, error)
	// send query !i<set-name>,1
	GetAsSetMembersRecrusive(set string) ([]uint32, error)
	// send query !gAS<asn> or !6AS<asn> depending on p
	GetRoutesByOrigin(p IPProto, asn uint32) ([]netip.Prefix, error)
	// send query !a<as-set-name> !a4<as-set-name> or !a6<as-set-name> depending on p
	GetRoutesBySet(p IPProto, set string) ([]netip.Prefix, error)

	io.Closer
}

type whois struct {
	proto.Client
}

func New(address string, timeout time.Duration) (Whois, error) {
	tcpConn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	dc := conn.NewDeadlineConn(tcpConn, timeout)
	return NewFromNetConn(dc), nil
}

func NewFromNetConn(c net.Conn) Whois {
	bc := conn.NewBufferedConn(c)
	return NewFromBufferedConn(bc)
}

func NewFromBufferedConn(c conn.BufferedConn) Whois {
	client := proto.New(c)
	return whois{
		Client: client,
	}
}

func (w whois) EnableMultiCommand() error {
	return w.Client.Write("!!")
}

func (w whois) SetIdleTimout(d time.Duration) error {
	_, err := w.Client.Query("!t" + fmt.Sprintf("%d", int(d.Seconds())))
	return err
}

func (w whois) GetVersion() (string, error) {
	return w.Client.Query("!v")
}

func (w whois) Close() error {
	return w.Client.Close()
}
