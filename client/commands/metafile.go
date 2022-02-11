package commands

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

// Struct that represents the data of the metafile
type MetafileData struct {
	Username       string
	Address        string
	FileExceptions map[string]struct{}
}

const METAFILE_NAME = "./.cvc"

var ErrMissingMetafile = errors.New("cannot open .cvc file")

// Tries to open the metafile to read it.
// Returns metafileData struct.
func ReadMetafileData() (*MetafileData, error) {
	file, err := openMetaFile()
	if err != nil {
		return nil, ErrMissingMetafile
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var username string
	var address string
	if scanner.Scan() && scanner.Err() == nil {
		username = scanner.Text()
	}

	if scanner.Scan() && scanner.Err() == nil {
		address = scanner.Text()
	}

	exceptions := make(map[string]struct{})
	for scanner.Scan() && scanner.Err() == nil {
		exceptions[scanner.Text()] = struct{}{}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &MetafileData{Username: username, Address: address, FileExceptions: exceptions}, nil
}

// Creates metafile and returns descriptor.
// Caller must close it.
func newMetaFile() (*os.File, error) {
	file, err := os.Create(METAFILE_NAME)
	if err != nil {
		return nil, fmt.Errorf("cannot create metafile: %w", err)
	}
	return file, nil
}

// Opens metafile and returns descriptor.
// Caller must close it.
func openMetaFile() (*os.File, error) {
	file, err := os.Open(METAFILE_NAME)
	if err != nil {
		return nil, fmt.Errorf("cannot open metafile: %w", err)
	}
	return file, nil
}

// Saves the metadata struct to a new metadata file
func (m *MetafileData) Save() error {
	file, err := newMetaFile()
	if err != nil {
		return err
	}
	defer file.Close()

	bufWriter := bufio.NewWriter(file)
	defer bufWriter.Flush()
	bufWriter.WriteString(m.Username)
	bufWriter.WriteRune('\n')
	bufWriter.WriteString(m.Address)
	bufWriter.WriteRune('\n')
	for files := range m.FileExceptions {
		bufWriter.WriteString(files)
		bufWriter.WriteRune('\n')
	}

	return nil
}
