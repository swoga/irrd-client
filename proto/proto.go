package proto

import (
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/swoga/irrd-client/conn"
	"github.com/swoga/irrd-client/expect"
)

type Client interface {
	io.Closer
	Query(query string) (string, error)
	Write(data string) error
}

type client struct {
	bufCon conn.BufferedConn
}

func New(c conn.BufferedConn) Client {
	return &client{c}
}

func (c client) Close() error {
	return c.bufCon.Close()
}

func (c *client) Write(data string) error {
	buffer := []byte(data + "\n")
	n, err := c.bufCon.Write(buffer)
	if err != nil {
		return err
	}
	if n != len(buffer) {
		return errors.New("incomplete write")
	}

	err = c.bufCon.Flush()
	if err != nil {
		return err
	}
	return nil
}

func (c *client) readTillNewline() (string, error) {
	buffer, err := c.bufCon.ReadBytes('\n')
	if err != nil {
		return "", err
	}
	// remove trailing \n
	buffer = buffer[:len(buffer)-1]
	return string(buffer), nil
}

func (c *client) readMsgA() (string, error) {
	lengthStr, err := c.readTillNewline()
	if err != nil {
		return "", err
	}
	length, err := strconv.Atoi(lengthStr)
	if err != nil {
		return "", err
	}
	// The length is the number of bytes in the response, including the newline immediately after the response content.
	length--

	buffer := make([]byte, length)
	n, err := io.ReadFull(c.bufCon, buffer)
	if err != nil {
		return "", err
	}
	if n != len(buffer) {
		return "", errors.New("incomplete read")
	}

	err = expect.String(c.bufCon, "\nC")
	if err != nil {
		return "", err
	}

	return string(buffer), nil
}

func (c *client) readMsgF() error {
	// read space
	err := expect.String(c.bufCon, " ")
	if err != nil {
		return err
	}
	msg, err := c.readTillNewline()
	if err != nil {
		return fmt.Errorf("error while reading error message: %v", err)
	}
	return fmt.Errorf("server returned error: %v", msg)
}

func (c *client) Query(query string) (string, error) {
	err := c.Write(query)
	if err != nil {
		return "", err
	}

	r, err := c.bufCon.ReadByte()
	if err != nil {
		return "", err
	}

	var result string

	switch string(r) {
	case "C":
		// query was valid, but no entries were found
	case "A":
		result, err = c.readMsgA()
	case "D":
		// query was valid, but the primary key queried for did not exist
	case "F":
		err = c.readMsgF()
	default:
		return "", fmt.Errorf("unknown response: %s", string(r))
	}

	if err != nil {
		return "", err
	}

	// handle trailing newline
	err = expect.String(c.bufCon, "\n")
	if err != nil {
		return "", err
	}

	return result, nil
}
