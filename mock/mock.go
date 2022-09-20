package mock

import (
	"bufio"
	"io"
	"net"
	"testing"

	"github.com/swoga/irrd-client/conn"
	"github.com/swoga/irrd-client/proto"
)

func SetupMockConn() (conn.BufferedConn, FakeServer) {
	serverConn, clientConn := net.Pipe()

	bc := conn.NewBufferedConn(clientConn)

	serverReader := bufio.NewReader(serverConn)
	server := FakeServer{reader: serverReader, writer: serverConn, Closer: serverConn}

	return bc, server
}

func SetupMockClient() (proto.Client, FakeServer) {
	bc, server := SetupMockConn()
	client := proto.New(bc)

	return client, server
}

type FakeServer struct {
	reader *bufio.Reader
	writer io.Writer
	io.Closer
}

func (f *FakeServer) Read(t *testing.T, want string) {
	has, err := f.reader.ReadString('\n')
	if err != nil {
		t.Fatal(err)
		return
	}
	if has != want {
		t.Fatalf("server got: %s, want: %s", has, want)
		return
	}
}

func (f *FakeServer) Write(t *testing.T, data string) {
	_, err := io.WriteString(f.writer, data)
	if err != nil {
		t.Fatal(err)
		return
	}
}
