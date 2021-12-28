package commands

import (
	"bufio"
	"errors"

	"github.com/NikolayDPaev/CentralisedVersionControl/client/fileIO"
)

type MetafileData struct {
	Username string
	Address  string
}

var ErrMissingMetafile = errors.New("cannot open .cvc file")

func ReadMetafileData() (*MetafileData, error) {
	file, err := fileIO.OpenMetaFile()
	if err != nil {
		return nil, ErrMissingMetafile
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	username := scanner.Text()
	scanner.Scan()
	address := scanner.Text()

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &MetafileData{Username: username, Address: address}, nil
}

func (m *MetafileData) Save() error {
	file, err := fileIO.NewMetaFile()
	if err != nil {
		return err
	}
	defer file.Close()

	bufWriter := bufio.NewWriter(file)
	defer bufWriter.Flush()
	bufWriter.WriteString(m.Username + "\n")
	bufWriter.WriteString(m.Address + "\n")

	return nil
}
