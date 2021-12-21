package commit

import (
	"bufio"
	"io"
)

type Metadata struct {
	message string
	creator string
	hash    string
}

func NewMetadataFromFile(reader io.Reader, hash string) (*Metadata, error) {
	scanner := bufio.NewScanner(reader)
	scanner.Scan()
	message := scanner.Text()
	scanner.Scan()
	creator := scanner.Text()

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &Metadata{message: message, creator: creator, hash: hash}, nil
}

func (m Metadata) String() string {
	return m.hash + " \"" + m.message + "\" " + m.creator
}
