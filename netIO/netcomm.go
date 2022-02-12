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
	DEF_CHUNK_SIZE   = 4096
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

func NewCommunication(chunkSize int, writer io.Writer, reader io.Reader) Communicator {
	return &NetCommunication{chunkSize, writer, reader}
}

// Writes an int with variable length to the io writer.
// Returns error if there is an IO error.
func (c *NetCommunication) SendVarInt(num int64) error {
	buf := make([]byte, INT_64_BYTE_LEN)
	bytes := binary.PutVarint(buf, num)

	if n, err := c.writer.Write(buf[:bytes]); err != nil || n != bytes {
		return fmt.Errorf("error sending var int %d: %w", num, err)
	}
	return nil
}

// Reads an int with variable length from the io reader.
// Returns error if there is an IO error.
func (c *NetCommunication) RecvVarInt() (int64, error) {
	num, err := ReadVarint(c.reader)
	if err != nil {
		return 0, fmt.Errorf("error reading var int: %w", err)
	}
	return num, nil
}

// Writes a string to the io writer.
// First sends a var int with its length then sends byte buffer with the string contains.
// Returns error if there is an IO error.
func (c *NetCommunication) SendString(str string) error {
	strBuf := []byte(str)
	err := c.SendVarInt(int64(len(strBuf)))
	if err != nil {
		return fmt.Errorf("could not send length of string: %w", err)
	}
	if n, err := c.writer.Write(strBuf); err != nil || n != len(strBuf) {
		return fmt.Errorf("error sending string %s: %w", str, err)
	}
	return nil
}

// Reads a string from the io reader.
// First reads a var int with the length then reads up to length bytes to a byte buffer.
// Returns error if there is an IO error.
func (c *NetCommunication) RecvString() (string, error) {
	len, err := c.RecvVarInt()
	if err != nil {
		return "", fmt.Errorf("could not recv length of string: %w", err)
	}

	bytes := make([]byte, len)
	n, err := c.reader.Read(bytes)

	if n != int(len) {
		return "", errors.New("length does not match")
	}
	if err != nil {
		return "", fmt.Errorf("could not recv string: %w", err)
	}

	return string(bytes), nil
}

// Writes a string slice to io writer.
// First sends a var int with its length, then sends the elements of the slice with sendString.
// Returns error if there is an IO error.
func (c *NetCommunication) SendStringSlice(slice []string) error {
	if err := c.SendVarInt(int64(len(slice))); err != nil {
		return fmt.Errorf("error sending string slice size: %w", err)
	}

	for _, str := range slice {
		if err := c.SendString(str); err != nil {
			return fmt.Errorf("error sending string slice: %w", err)
		}
	}
	return nil
}

// Reads a string slice from the io reader.
// First reads a var int with its length, then reads the strings with recvString.
// Returns error if there is an IO error.
func (c *NetCommunication) RecvStringSlice() ([]string, error) {
	len, err := c.RecvVarInt()
	if err != nil {
		return nil, fmt.Errorf("error receiving string slice size: %w", err)
	}

	slice := make([]string, len)
	for i := 0; i < int(len); i++ {
		slice[i], err = c.RecvString()
		if err != nil {
			return nil, fmt.Errorf("error receiving string slice: %w", err)
		}
	}
	return slice, nil
}

// Reads up to fileLength bytes from the fileReader and writes them to the io writer.
// It is done on chunks with the specified chunkSize.
// Returns error if there is an IO error.
func (c *NetCommunication) SendFileData(fileReader io.Reader, fileLength int64) error {
	err := c.SendVarInt(fileLength)
	if err != nil {
		return fmt.Errorf("error sending file length: %w", err)
	}

	buf := make([]byte, c.chunkSize)

	bytesRead, readErr := fileReader.Read(buf)
	_, sendErr := c.writer.Write(buf[:bytesRead])

	for bytesRead > 0 && readErr == nil && sendErr == nil {
		bytesRead, readErr = fileReader.Read(buf)
		_, sendErr = c.writer.Write(buf[:bytesRead])
	}
	if readErr != nil && !errors.Is(readErr, io.EOF) {
		return fmt.Errorf("error reading file: %w", readErr)
	}
	if sendErr != nil {
		return fmt.Errorf("error sending file: %w", sendErr)
	}

	return nil
}

// Reads a number from the io reader and then reads up to the this number
// from the reader and writes it to the provided fileWriter.
// It is done on chunks with the specified chunkSize.
// Returns error if there is an IO error.
func (c *NetCommunication) RecvFileData(fileWriter io.Writer) error {
	remaining, err := c.RecvVarInt()
	if err != nil {
		return fmt.Errorf("error receiving file length: %w", err)
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
		return fmt.Errorf("error receiving file: %w", readErr)
	}
	if sendErr != nil {
		return fmt.Errorf("error writing file: %w", sendErr)
	}

	return nil
}
