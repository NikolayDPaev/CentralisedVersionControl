package commit

import (
	"fmt"
	"io"
	"strings"

	"github.com/NikolayDPaev/CentralisedVersionControl/server/netIO"
)

type Commit struct {
	id      string
	message string
	creator string
	tree    string
}

func ReadCommitData(reader io.Reader) (string, string, error) {
	message, err := netIO.ReceiveString(reader)
	if err != nil {
		return "", "", err
	}

	creator, err := netIO.ReceiveString(reader)
	if err != nil {
		return "", "", err
	}
	return message, creator, nil
}

func ReadCommit(reader io.Reader) (*Commit, error) {
	id, err := netIO.ReceiveString(reader)
	if err != nil {
		return nil, fmt.Errorf("cannot read id of commit:\n%w", err)
	}

	message, creator, err := ReadCommitData(reader)
	if err != nil {
		return nil, fmt.Errorf("cannot read metadata of commit:\n%w", err)
	}

	tree, err := netIO.ReceiveString(reader)
	if err != nil {
		return nil, fmt.Errorf("cannot read tree string of commit:\n%w", err)
	}

	return &Commit{id, message, creator, tree}, nil
}

func (c *Commit) Id() string {
	return c.id
}

func (c *Commit) Write(writer io.Writer) error {
	err := netIO.SendString(c.message, writer)
	if err != nil {
		return fmt.Errorf("cannot send commit message:\n%w", err)
	}

	err = netIO.SendString(c.creator, writer)
	if err != nil {
		return fmt.Errorf("cannot send commit creator:\n%w", err)
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
	return c.id + " \"" + c.message + "\" " + c.creator
}
