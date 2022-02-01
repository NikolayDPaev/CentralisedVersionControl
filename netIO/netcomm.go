package netio

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

const (
	UINT_32_BYTE_LEN = 4
	INT_64_BYTE_LEN  = 8
)

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

type NetCommunication struct {
	chunkSize int
	writer    io.Writer
	reader    io.Reader
}

func NewCommunicator(chunkSize int, writer io.Writer, reader io.Reader) Communicator {
	return &NetCommunication{chunkSize, writer, reader}
}

func (c *NetCommunication) SendVarInt(num int64) error {
	buf := make([]byte, INT_64_BYTE_LEN)
	bytes := binary.PutVarint(buf, num)

	if n, err := c.writer.Write(buf[:bytes]); err != nil || n != bytes {
		return fmt.Errorf("error sending var int %d:\n%w", num, err)
	}
	return nil
}

func (c *NetCommunication) RecvVarInt() (int64, error) {
	num, err := ReadVarint(c.reader)
	if err != nil {
		return 0, fmt.Errorf("error reading var int:\n%w", err)
	}
	return num, nil
}

func (c *NetCommunication) SendString(str string) error {
	strBuf := []byte(str)
	err := c.SendVarInt(int64(len(strBuf)))
	if err != nil {
		return fmt.Errorf("could not send length of string:\n%w", err)
	}
	if n, err := c.writer.Write(strBuf); err != nil || n != len(strBuf) {
		return fmt.Errorf("error sending string %s:\n%w", str, err)
	}
	return nil
}

func (c *NetCommunication) RecvString() (string, error) {
	len, err := c.RecvVarInt()
	if err != nil {
		return "", fmt.Errorf("could not recv length of string:\n%w", err)
	}

	bytes := make([]byte, len)
	n, err := c.reader.Read(bytes)

	if n != int(len) {
		return "", errors.New("length does not match")
	}
	if err != nil {
		return "", fmt.Errorf("could not recv string:\n%w", err)
	}

	return string(bytes), nil
}

func (c *NetCommunication) SendStringSlice(slice []string) error {
	if err := c.SendVarInt(int64(len(slice))); err != nil {
		return fmt.Errorf("error sending string slice size:\n%w", err)
	}

	for _, str := range slice {
		if err := c.SendString(str); err != nil {
			return fmt.Errorf("error sending string slice:\n%w", err)
		}
	}
	return nil
}

func (c *NetCommunication) RecvStringSlice() ([]string, error) {
	len, err := c.RecvVarInt()
	if err != nil {
		return nil, fmt.Errorf("error receiving string slice size:\n%w", err)
	}

	slice := make([]string, len)
	for i := 0; i < int(len); i++ {
		slice[i], err = c.RecvString()
		if err != nil {
			return nil, fmt.Errorf("error receiving string slice:\n%w", err)
		}
	}
	return slice, nil
}

func (c *NetCommunication) SendFileData(fileReader io.Reader, fileLength int64) error {
	err := c.SendVarInt(fileLength)
	if err != nil {
		return fmt.Errorf("error sending file length:\n%w", err)
	}

	buf := make([]byte, c.chunkSize)

	bytesRead, readErr := fileReader.Read(buf)
	_, sendErr := c.writer.Write(buf[:bytesRead])

	for bytesRead > 0 && readErr == nil && sendErr == nil {
		bytesRead, readErr = fileReader.Read(buf)
		_, sendErr = c.writer.Write(buf[:bytesRead])
	}
	if readErr != nil && !errors.Is(readErr, io.EOF) {
		return fmt.Errorf("error reading file:\n%w", readErr)
	}
	if sendErr != nil {
		return fmt.Errorf("error sending file:\n%w", sendErr)
	}

	return nil
}

func (c *NetCommunication) RecvFileData(fileWriter io.Writer) error {
	remaining, err := c.RecvVarInt()
	if err != nil {
		return fmt.Errorf("error receiving file length:\n%w", err)
	}
	buf := make([]byte, c.chunkSize)

	var readErr error
	var sendErr error
	var bytesRead int

	for remaining > 0 && readErr == nil && sendErr == nil {
		bytesRead, readErr = c.reader.Read(buf[:min(remaining, int64(c.chunkSize))])
		remaining -= int64(bytesRead)
		_, sendErr = fileWriter.Write(buf[:bytesRead])
	}
	if readErr != nil {
		return fmt.Errorf("error receiving file:\n%w", readErr)
	}
	if sendErr != nil {
		return fmt.Errorf("error writing file:\n%w", sendErr)
	}

	return nil
}
