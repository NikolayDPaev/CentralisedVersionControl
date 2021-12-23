package netIO

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

const (
	CHUNK_SIZE       = 4096
	UINT_32_BYTE_LEN = 4
	INT_64_BYTE_LEN  = 8
)

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func SendVarInt(num int64, writer io.Writer) error {
	buf := make([]byte, INT_64_BYTE_LEN)
	bytes := binary.PutVarint(buf, num)

	if n, err := writer.Write(buf[:bytes]); err != nil || n != bytes {
		return fmt.Errorf("error sending var int %d: %w", num, err)
	}
	return nil
}

func ReceiveVarInt(reader io.Reader) (int64, error) {
	byteReader := io.ByteReader(bufio.NewReader(reader))

	num, err := binary.ReadVarint(byteReader)
	if err == nil {
		return 0, fmt.Errorf("error reading var int: %w", err)
	}
	return num, nil
}

func SendUint32(num uint32, writer io.Writer) error {
	bytes := make([]byte, UINT_32_BYTE_LEN)
	binary.LittleEndian.PutUint32(bytes, num)
	n, err := writer.Write(bytes)

	if n != UINT_32_BYTE_LEN {
		return errors.New("could not send number")
	}
	if err != nil {
		return fmt.Errorf("could not send number: %w", err)
	}
	return nil
}

func SendString(str string, writer io.Writer) error {
	err := SendUint32(uint32(len(str)), writer)
	if err != nil {
		return fmt.Errorf("could not send length of string: %w", err)
	}
	writer.Write([]byte(str))
	return nil
}

func ReceiveUint32(reader io.Reader) (uint32, error) {
	bytes := make([]byte, UINT_32_BYTE_LEN)
	n, err := reader.Read(bytes)

	if n != UINT_32_BYTE_LEN {
		return 0, errors.New("length does not match")
	}
	if err != nil {
		return 0, fmt.Errorf("could not receive number: %w", err)
	}
	return binary.BigEndian.Uint32(bytes), nil
}

func SendFileData(fileReader io.Reader, fileLength int64, netWriter io.Writer) error {
	SendVarInt(fileLength, netWriter)

	bufFileReader := bufio.NewReader(fileReader)
	bufNetWriter := bufio.NewWriter(netWriter)
	defer bufNetWriter.Flush()

	buf := make([]byte, CHUNK_SIZE)

	bytesRead, readErr := bufFileReader.Read(buf)
	_, sendErr := bufNetWriter.Write(buf[:bytesRead])

	for bytesRead > 0 && readErr == nil && sendErr == nil {
		bytesRead, readErr = bufFileReader.Read(buf)
		_, sendErr = bufNetWriter.Write(buf[:bytesRead])
	}
	if readErr != nil {
		return fmt.Errorf("error reading file: %w", readErr)
	}
	if sendErr != nil {
		return fmt.Errorf("error sending file: %w", sendErr)
	}

	return nil
}

func ReceiveFileData(netReader io.Reader, fileWriter io.Writer) error {
	remaining, err := ReceiveVarInt(netReader)
	if err != nil {
		return fmt.Errorf("error receiving file length: %w", err)
	}
	bufNetReader := bufio.NewReader(netReader)
	bufFileWriter := bufio.NewWriter(fileWriter)
	buf := make([]byte, CHUNK_SIZE)

	var readErr error
	var sendErr error
	var bytesRead int

	for remaining > 0 && readErr == nil && sendErr == nil {
		bytesRead, readErr = bufNetReader.Read(buf[:min(remaining, CHUNK_SIZE)])
		remaining -= int64(bytesRead)
		_, sendErr = bufFileWriter.Write(buf[:bytesRead])
	}
	if readErr != nil {
		return fmt.Errorf("error receiving file: %w", readErr)
	}
	if sendErr != nil {
		return fmt.Errorf("error writing file: %w", sendErr)
	}

	return nil
}
