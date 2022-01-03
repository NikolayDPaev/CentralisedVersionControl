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
		return fmt.Errorf("error sending var int %d:\n%w", num, err)
	}
	return nil
}

func ReceiveVarInt(reader io.Reader) (int64, error) {
	num, err := ReadVarint(reader)
	if err != nil {
		return 0, fmt.Errorf("error reading var int:\n%w", err)
	}
	return num, nil
}

func SendString(str string, writer io.Writer) error {
	strBuf := []byte(str)
	err := SendVarInt(int64(len(strBuf)), writer)
	if err != nil {
		return fmt.Errorf("could not send length of string:\n%w", err)
	}
	if n, err := writer.Write(strBuf); err != nil || n != len(strBuf) {
		return fmt.Errorf("error sending string %s:\n%w", str, err)
	}
	return nil
}

func ReceiveString(reader io.Reader) (string, error) {
	len, err := ReceiveVarInt(reader)
	if err != nil {
		return "", fmt.Errorf("could not receive length of string:\n%w", err)
	}

	bytes := make([]byte, len)
	n, err := reader.Read(bytes)

	if n != int(len) {
		return "", errors.New("length does not match")
	}
	if err != nil {
		return "", fmt.Errorf("could not receive string:\n%w", err)
	}

	return string(bytes), nil
}

func SendStringSlice(slice []string, writer io.Writer) error {
	if err := SendVarInt(int64(len(slice)), writer); err != nil {
		return fmt.Errorf("error sending string slice size:\n%w", err)
	}

	for _, str := range slice {
		if err := SendString(str, writer); err != nil {
			return fmt.Errorf("error sending string slice:\n%w", err)
		}
	}
	return nil
}

func ReceiveStringSlice(reader io.Reader) ([]string, error) {
	len, err := ReceiveVarInt(reader)
	if err != nil {
		return nil, fmt.Errorf("error receiving string slice size:\n%w", err)
	}

	slice := make([]string, len)
	for i := 0; i < int(len); i++ {
		slice[i], err = ReceiveString(reader)
		if err != nil {
			return nil, fmt.Errorf("error receiving string slice:\n%w", err)
		}
	}
	return slice, nil
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
		return fmt.Errorf("error reading file:\n%w", readErr)
	}
	if sendErr != nil {
		return fmt.Errorf("error sending file:\n%w", sendErr)
	}

	return nil
}

func ReceiveFileData(netReader io.Reader, fileWriter io.Writer) error {
	remaining, err := ReceiveVarInt(netReader)
	if err != nil {
		return fmt.Errorf("error receiving file length:\n%w", err)
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
		return fmt.Errorf("error receiving file:\n%w", readErr)
	}
	if sendErr != nil {
		return fmt.Errorf("error writing file:\n%w", sendErr)
	}

	return nil
}
