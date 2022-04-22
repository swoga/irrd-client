package conn

import (
	"bufio"
	"net"
)

type BufferedConn interface {
	net.Conn
	ReadByte() (byte, error)
	ReadBytes(delim byte) ([]byte, error)
	Flush() error
}

func NewBufferedConn(conn net.Conn) BufferedConn {
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	return bufferedConn{Conn: conn, reader: reader, writer: writer}
}

type bufferedConn struct {
	net.Conn
	reader *bufio.Reader
	writer *bufio.Writer
}

func (c bufferedConn) Read(b []byte) (int, error) {
	return c.reader.Read(b)
}

func (c bufferedConn) ReadByte() (byte, error) {
	return c.reader.ReadByte()
}

func (c bufferedConn) ReadBytes(delim byte) ([]byte, error) {
	return c.reader.ReadBytes(delim)
}

func (c bufferedConn) Write(b []byte) (int, error) {
	return c.writer.Write(b)
}

func (c bufferedConn) Flush() error {
	return c.writer.Flush()
}
