package fileio

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const (
	CHUNK_SIZE = 4096
)

func (l *Localfiles) CompressToTempFile(source string) (*os.File, error) {
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

func (l *Localfiles) createDirectoriesInPath(dest string) error {
	if err := os.MkdirAll(filepath.Dir(dest), 0777); err != nil {
		return fmt.Errorf("cannot create path: %w", err)
	}
	return nil
}

func (l *Localfiles) DecompressFile(dest string, sFile *os.File) error {
	if err := l.createDirectoriesInPath(dest); err != nil {
		return fmt.Errorf("error creating file directory: %w", err)
	}
	dFile, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
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
