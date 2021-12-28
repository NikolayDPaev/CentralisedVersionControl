package commit

import (
	"fmt"
	"io"

	"github.com/NikolayDPaev/CentralisedVersionControl/server/netIO"
)

type Metadata struct {
	id      string
	message string
	creator string
}

func ReadMetadata(reader io.Reader, id string) (*Metadata, error) {
	message, err := netIO.ReceiveString(reader)
	if err != nil {
		return nil, err
	}

	creator, err := netIO.ReceiveString(reader)
	if err != nil {
		return nil, err
	}
	return &Metadata{message: message, creator: creator, id: id}, nil
}

func (m Metadata) String() string {
	return m.id + " \"" + m.message + "\" " + m.creator
}

func (m *Metadata) Write(writer io.Writer) error {
	err := netIO.SendString(m.message, writer)
	if err != nil {
		return fmt.Errorf("cannot send commit message: %w", err)
	}

	err = netIO.SendString(m.creator, writer)
	if err != nil {
		return fmt.Errorf("cannot send commit creator: %w", err)
	}
	return nil
}
