package metadata

import (
	"bufio"
	"errors"
)

var ErrMissingMetafile = errors.New("cannot open .cvc file")

// MetafileData is a struct that represents the data of the metafile
type MetafileData struct {
	Username     string
	Address      string
	IgnoredFiles map[string]struct{}
}

// ReadMetafileData tries to open the metafile to read it.
// Returns metafileData struct.
func ReadMetafileData(metafileName string) (*MetafileData, error) {
	file, err := openMetaFile(metafileName)
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

	ignored := make(map[string]struct{})
	for scanner.Scan() && scanner.Err() == nil {
		ignored[scanner.Text()] = struct{}{}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &MetafileData{Username: username, Address: address, IgnoredFiles: ignored}, nil
}

// Save saves the metadata struct to a new metafile
func Save(data *MetafileData, metafileName string) error {
	file, err := newMetaFile(metafileName)
	if err != nil {
		return err
	}
	defer file.Close()

	bufWriter := bufio.NewWriter(file)
	defer bufWriter.Flush()
	bufWriter.WriteString(data.Username)
	bufWriter.WriteRune('\n')
	bufWriter.WriteString(data.Address)
	bufWriter.WriteRune('\n')
	for files := range data.IgnoredFiles {
		bufWriter.WriteString(files)
		bufWriter.WriteRune('\n')
	}

	return nil
}
