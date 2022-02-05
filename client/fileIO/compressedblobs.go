package fileio

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/NikolayDPaev/CentralisedVersionControl/netio"
)

const (
	CHUNK_SIZE = 4096
)

func (l *Localfiles) compressToTempFile(source string) (*os.File, error) {
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

func (l *Localfiles) decompressFile(dest string, sFile *os.File) error {
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

func (l *Localfiles) ReceiveBlob(filepath string, comm netio.Communicator) error {
	tmp, err := os.CreateTemp("", "blobTmp")
	if err != nil {
		return fmt.Errorf("error creating tmpBlob:\n%w", err)
	}
	defer tmp.Close()

	if err := comm.RecvFileData(tmp); err != nil {
		return fmt.Errorf("error receiving blob:\n%w", err)
	}

	if err := l.decompressFile(filepath, tmp); err != nil {
		return fmt.Errorf("error decompressing blob:\n%w", err)
	}
	return nil
}

func (l *Localfiles) SendBlob(filepath string, comm netio.Communicator) error {
	tmpFile, err := l.compressToTempFile(filepath)
	if err != nil {
		return fmt.Errorf("error compressing file %s:\n%w", filepath, err)
	}
	defer tmpFile.Close()

	stat, err := tmpFile.Stat()
	if err != nil {
		return fmt.Errorf("error getting blobTmp size:\n%w", err)
	}

	err = comm.SendFileData(tmpFile, stat.Size())
	if err != nil {
		return fmt.Errorf("error sending blob:\n%w", err)
	}
	return nil
}
