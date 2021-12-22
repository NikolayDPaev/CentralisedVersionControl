package clienthandler

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

func sendUint32(num uint32, writer io.Writer) error {
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, num)
	n, err := writer.Write(bytes)

	if n != 4 {
		return errors.New("could not send number")
	}
	if err != nil {
		return fmt.Errorf("could not send number: %w", err)
	}
	return nil
}

func sendString(str string, writer io.Writer) error {
	err := sendUint32(uint32(len(str)), writer)
	if err != nil {
		return fmt.Errorf("could not send length of string: %w", err)
	}
	writer.Write([]byte(str))
	return nil
}

func receiveUint32(reader io.Reader) (uint32, error) {
	bytes := make([]byte, 4)
	n, err := reader.Read(bytes)

	if n != 4 {
		return 0, errors.New("length does not match")
	}
	if err != nil {
		return 0, fmt.Errorf("could not receive number: %w", err)
	}
	return binary.BigEndian.Uint32(bytes), nil
}
