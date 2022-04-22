package mock

import (
	"bufio"
	"io"
	"net"

	"github.com/swoga/irrd-client/asynctest"
	"github.com/swoga/irrd-client/conn"
	"github.com/swoga/irrd-client/proto"
)

func SetupMockConn() (conn.BufferedConn, FakeServer) {
	serverConn, clientConn := net.Pipe()

	bc := conn.NewBufferedConn(clientConn)

	serverReader := bufio.NewReader(serverConn)
	server := FakeServer{reader: serverReader, writer: serverConn}

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
}

func (f *FakeServer) Read(t asynctest.T, want string) {
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

func (f *FakeServer) Write(t asynctest.T, data string) {
	_, err := io.WriteString(f.writer, data)
	if err != nil {
		t.Fatal(err)
		return
	}
}
