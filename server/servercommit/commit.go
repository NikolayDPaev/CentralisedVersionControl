package servercommit

import (
	"fmt"
	"strings"

	"github.com/NikolayDPaev/CentralisedVersionControl/netIO"
)

type Commit struct {
	id      string
	message string
	creator string
	tree    string
}

func NewCommit(id, message, creator, tree string) *Commit {
	return &Commit{id, message, creator, tree}
}

func ReadCommitData(comm netIO.Communicator) (string, string, error) {
	message, err := comm.ReceiveString()
	if err != nil {
		return "", "", err
	}

	creator, err := comm.ReceiveString()
	if err != nil {
		return "", "", err
	}
	return message, creator, nil
}

func ReadCommit(id string, comm netIO.Communicator) (*Commit, error) {
	message, creator, err := ReadCommitData(comm)
	if err != nil {
		return nil, fmt.Errorf("cannot read metadata of commit:\n%w", err)
	}

	tree, err := comm.ReceiveString()
	if err != nil {
		return nil, fmt.Errorf("cannot read tree string of commit:\n%w", err)
	}

	return &Commit{id, message, creator, tree}, nil
}

func (c *Commit) Id() string {
	return c.id
}

func (c *Commit) Write(comm netIO.Communicator) error {
	err := comm.SendString(c.message)
	if err != nil {
		return fmt.Errorf("cannot send commit message:\n%w", err)
	}

	err = comm.SendString(c.creator)
	if err != nil {
		return fmt.Errorf("cannot send commit creator:\n%w", err)
	}

	if err := comm.SendString(c.tree); err != nil {
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
