package commit

import (
	"bufio"
	"fmt"
	"io"
)

type Metadata struct {
	id      string
	message string
	creator string
}

func ReadMetadata(reader io.Reader, id string) (*Metadata, error) {
	scanner := bufio.NewScanner(reader)
	scanner.Scan()
	message := scanner.Text()
	scanner.Scan()
	creator := scanner.Text()

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &Metadata{message: message, creator: creator, id: id}, nil
}

func (m Metadata) String() string {
	return m.id + " \"" + m.message + "\" " + m.creator
}

func (m *Metadata) Write(writer io.Writer) error {
	bufWriter := bufio.NewWriter(writer)

	_, err := bufWriter.WriteString(m.message + "\n")
	if err != nil {
		return fmt.Errorf("cannot write metadata: %w", err)
	}

	_, err = bufWriter.WriteString(m.creator + "\n")
	if err != nil {
		return fmt.Errorf("cannot write metadata: %w", err)
	}
	return nil
}
