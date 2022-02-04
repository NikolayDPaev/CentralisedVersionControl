package servercommit

import (
	"fmt"
	"strings"

	"github.com/NikolayDPaev/CentralisedVersionControl/netio"
)

type Commit struct {
	Id      string
	Message string
	Creator string
	Tree    string
}

func ReadCommitData(comm netio.Communicator) (string, string, error) {
	message, err := comm.RecvString()
	if err != nil {
		return "", "", err
	}

	creator, err := comm.RecvString()
	if err != nil {
		return "", "", err
	}
	return message, creator, nil
}

func ReadCommit(id string, comm netio.Communicator) (*Commit, error) {
	message, creator, err := ReadCommitData(comm)
	if err != nil {
		return nil, fmt.Errorf("cannot read metadata of commit:\n%w", err)
	}

	tree, err := comm.RecvString()
	if err != nil {
		return nil, fmt.Errorf("cannot read tree string of commit:\n%w", err)
	}

	return &Commit{id, message, creator, tree}, nil
}

func (c *Commit) WriteData(comm netio.Communicator) error {
	err := comm.SendString(c.Message)
	if err != nil {
		return fmt.Errorf("cannot send commit message:\n%w", err)
	}

	err = comm.SendString(c.Creator)
	if err != nil {
		return fmt.Errorf("cannot send commit creator:\n%w", err)
	}

	if err := comm.SendString(c.Tree); err != nil {
		return fmt.Errorf("cannot write commit tree:\n%w", err)
	}

	return nil
}

func (c *Commit) ExtractBlobIds() []string { // regex ????
	lines := strings.Split(c.Tree, "\n")

	blobIds := make([]string, len(lines))
	for i, line := range lines {
		blobIds[i] = strings.Split(line, " ")[0]
	}

	return blobIds
}

func (c Commit) String() string {
	return c.Id + " \"" + c.Message + "\" " + c.Creator
}
