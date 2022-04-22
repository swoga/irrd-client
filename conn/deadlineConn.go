package conn

import (
	"net"
	"time"
)

type deadlineConn struct {
	net.Conn
	timeout time.Duration
}

func NewDeadlineConn(conn net.Conn, timeout time.Duration) net.Conn {
	return deadlineConn{conn, timeout}
}

func (c deadlineConn) Read(b []byte) (int, error) {
	c.SetReadDeadline(time.Now().Add(c.timeout))
	return c.Conn.Read(b)
}

func (c deadlineConn) Write(b []byte) (int, error) {
	c.SetWriteDeadline(time.Now().Add(c.timeout))
	return c.Conn.Write(b)
}
