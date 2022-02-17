package fileio

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/NikolayDPaev/CentralisedVersionControl/netio"
)

// compressToTempFile compresses the provided file to new temporary file.
// Returns descriptor of the temp file and it is caller responsibiliy to close it.
// Uses gzip for the compresion.
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
	defer gzipWriter.Close()

	if _, err := io.Copy(gzipWriter, sFile); err != nil {
		return nil, fmt.Errorf("error while writing to gzipWriter the file file: %w", err)
	}

	gzipWriter.Close()
	dFile.Seek(0, io.SeekStart)
	return dFile, nil
}

// createDirectoriesInPath creates the directories leading to the file in the filepath.
func (l *Localfiles) createDirectoriesInPath(dest string) error {
	if err := os.MkdirAll(filepath.Dir(dest), 0777); err != nil {
		return fmt.Errorf("cannot create path: %w", err)
	}
	return nil
}

// decompressFile decompresses the provided file descriptor to a new file with the provided path.
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

	if _, err := io.Copy(dFile, gzipReader); err != nil {
		return fmt.Errorf("error while writing to gzipWriter the file file: %w", err)
	}

	return nil
}

// ReceiveBlob encapsulates the logic behind receiving blob.
// Receives it to a temporary file then decompresses it to a new file
//  with the specified path.
func (l *Localfiles) ReceiveBlob(filepath string, comm netio.Communicator) error {
	tmp, err := os.CreateTemp("", "blobTmp")
	if err != nil {
		return fmt.Errorf("error creating tmpBlob: %w", err)
	}
	defer tmp.Close()

	if err := comm.RecvFileData(tmp); err != nil {
		return fmt.Errorf("error receiving blob: %w", err)
	}

	tmp.Seek(0, io.SeekStart)

	if err := l.decompressFile(filepath, tmp); err != nil {
		return fmt.Errorf("error decompressing blob: %w", err)
	}
	return nil
}

// SendBlob encapsulates the logic behind sending blob.
// Compresses the file with the provided path to a temporary file
// then proceeds to send it via the Communicator interface.
func (l *Localfiles) SendBlob(filepath string, comm netio.Communicator) error {
	tmpFile, err := l.compressToTempFile(filepath)
	if err != nil {
		return fmt.Errorf("error compressing file %s: %w", filepath, err)
	}
	defer tmpFile.Close()

	stat, err := tmpFile.Stat()
	if err != nil {
		return fmt.Errorf("error getting blobTmp size: %w", err)
	}

	err = comm.SendFileData(tmpFile, stat.Size())
	if err != nil {
		return fmt.Errorf("error sending blob: %w", err)
	}
	return nil
}
