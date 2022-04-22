package expect

import (
	"errors"
	"io"
)

func Bytes(r io.Reader, value []byte) error {
	buffer := make([]byte, len(value))
	n, err := io.ReadFull(r, buffer)
	if err != nil {
		return err
	}
	if n != len(value) {
		return errors.New("incomplete read")
	}
	for i := range value {
		if value[i] != buffer[i] {
			return errors.New("unexpected byte")
		}
	}
	return nil
}

func String(r io.Reader, value string) error {
	return Bytes(r, []byte(value))
}
