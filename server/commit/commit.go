package commit

import (
	"fmt"
	"io"
	"strings"

	"github.com/NikolayDPaev/CentralisedVersionControl/server/netIO"
)

type Commit struct {
	metadata Metadata
	tree     string
}

func ReadCommit(reader io.Reader) (*Commit, error) {
	id, err := netIO.ReceiveString(reader)
	if err != nil {
		return nil, fmt.Errorf("cannot read id of commit:\n%w", err)
	}

	metadata, err := ReadMetadata(reader, id)
	if err != nil {
		return nil, fmt.Errorf("cannot read metadata of commit:\n%w", err)
	}

	tree, err := netIO.ReceiveString(reader)
	if err != nil {
		return nil, fmt.Errorf("cannot read tree string of commit:\n%w", err)
	}

	return &Commit{metadata: *metadata, tree: tree}, nil
}

func (c *Commit) Id() string {
	return c.metadata.id
}

func (c *Commit) Write(writer io.Writer) error {
	if err := c.metadata.Write(writer); err != nil {
		return fmt.Errorf("cannot write commit metadata:\n%w", err)
	}

	if err := netIO.SendString(c.tree, writer); err != nil {
		return fmt.Errorf("cannot write commit tree:\n%w", err)
	}

	return nil
}

func (c *Commit) ExtractBlobIds() []string { // regex ????
	lines := strings.Split(c.tree, "\n")

	blobIds := make([]string, len(lines))
	for i, line := range lines {
		blobIds[i] = strings.Split(line, " ")[0]
	}

	return blobIds
}

func (c Commit) String() string {
	return c.metadata.String()
}
