package commands

import (
	"bufio"
	"errors"

	"github.com/NikolayDPaev/CentralisedVersionControl/client/fileio"
)

type MetafileData struct {
	Username       string
	Address        string
	FileExceptions map[string]struct{}
}

var ErrMissingMetafile = errors.New("cannot open .cvc file")

func ReadMetafileData() (*MetafileData, error) {
	file, err := fileio.OpenMetaFile()
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

func (m *MetafileData) Save() error {
	file, err := fileio.NewMetaFile()
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
