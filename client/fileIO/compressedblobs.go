package fileIO

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
)

const (
	CHUNK_SIZE = 4096
)

func CompressToTempFile(source string) (*os.File, error) {
	sFile, err := os.Open(source)
	if err != nil {
		return nil, fmt.Errorf("error opening source file: %w", err)
	}
	defer sFile.Close()

	dFile, err := os.CreateTemp("", "blobTmp")
	if err != nil {
		return nil, fmt.Errorf("error creating temp file: %w", err)
	}

	gzipWriter := gzip.NewWriter(dFile)
	bytes := make([]byte, CHUNK_SIZE)

	n, errR := sFile.Read(bytes)
	_, errW := gzipWriter.Write(bytes[:n])
	for n > 0 && errR == nil && errW == nil {
		n, errR = sFile.Read(bytes)
		_, errW = gzipWriter.Write(bytes[:n])
	}
	if err != nil {
		return nil, fmt.Errorf("error compressing chunks of file: %w", err)
	}

	gzipWriter.Close()
	dFile.Seek(0, io.SeekStart)
	return dFile, nil
}

func DecompressFile(dest string, sFile *os.File) error {
	dFile, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("error creating decompressed file: %w", err)
	}
	defer sFile.Close()

	gzipReader, err := gzip.NewReader(sFile)
	if err != nil {
		return fmt.Errorf("error decompressing file: %w", err)
	}
	defer gzipReader.Close()
	bytes := make([]byte, CHUNK_SIZE)

	n, errR := gzipReader.Read(bytes)
	_, errW := dFile.Write(bytes[:n])
	for n > 0 && errR == nil && errW == nil {
		n, errR = gzipReader.Read(bytes)
		_, errW = dFile.Write(bytes[:n])
	}

	if err != nil {
		return fmt.Errorf("error decompressing chunks of file: %w", err)
	}

	return nil
}
