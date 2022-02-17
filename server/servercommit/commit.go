// Package that defines the Commit structure and methods for it serialization/deserialization
package servercommit

import (
	"fmt"
	"strings"

	"github.com/NikolayDPaev/CentralisedVersionControl/netio"
)

// Commit is a record that represents the server view of a commit.
type Commit struct {
	Id      string
	Message string
	Creator string
	Tree    string
}

// NewCommitFrom deserializes a new commit from the provided communicator.
// Returns an error if any of the receive operations fails.
func NewCommitFrom(id string, comm netio.Communicator) (*Commit, error) {
	message, err := comm.RecvString()
	if err != nil {
		return nil, err
	}

	creator, err := comm.RecvString()
	if err != nil {
		return nil, err
	}

	tree, err := comm.RecvString()
	if err != nil {
		return nil, fmt.Errorf("cannot read tree string of commit: %w", err)
	}

	return &Commit{id, message, creator, tree}, nil
}

// WriteTo serializes the commit to the provided communicator.
// Returns an error if any of the send operations fails.
func (c *Commit) WriteTo(comm netio.Communicator) error {
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

// ExtractBlobIds returns slice with the blob ids in the commit.
func (c *Commit) ExtractBlobIds() []string {
	lines := strings.Split(c.Tree, "\n")

	blobIds := make([]string, len(lines))
	for i, line := range lines {
		words := strings.Split(line, " ")
		if len(words) > 0 {
			blobIds[i] = words[0]
		}
	}

	return blobIds
}

// String returns string representation of the commit.
func (c Commit) String() string {
	return c.Id + " \"" + c.Message + "\" " + c.Creator
}
