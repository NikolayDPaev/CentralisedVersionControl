package commands

import (
	"bufio"
	"fmt"

	"github.com/NikolayDPaev/CentralisedVersionControl/client/fileIO"
)

type MetafileData struct {
	Username string
	Address  string
}

func ReadMetafileData() (*MetafileData, error) {
	file, err := fileIO.OpenMetaFile()
	if err != nil {
		fmt.Println("Cannot fing .cvc file. Please run command csv init")
		return nil, err
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
