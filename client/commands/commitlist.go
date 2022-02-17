package commands

import (
	"fmt"

	"github.com/NikolayDPaev/CentralisedVersionControl/netio"
)

// Commitlist implements the list commits operation
type Commitlist struct {
	comm   netio.Communicator
	opcode int
}

func NewCommitList(comm netio.Communicator, opcode int) *Commitlist {
	return &Commitlist{comm, opcode}
}

// GetCommitList requsts a commit list from the server
// Sends opcode of the operation and receives a slice from strings
// representing the commits on the server.
func (c *Commitlist) GetCommitList() ([]string, error) {
	if err := c.comm.SendVarInt(int64(c.opcode)); err != nil {
		return nil, fmt.Errorf("cannot send opcode: %w", err)
	}

	commitList, err := c.comm.RecvStringSlice()
	if err != nil {
		return nil, fmt.Errorf("error receiving commit list: %w", err)
	}
	return commitList, nil
}
