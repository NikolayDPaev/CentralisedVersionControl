package commands

import (
	"fmt"

	"github.com/NikolayDPaev/CentralisedVersionControl/netio"
)

type Commitlist struct {
	comm   netio.Communicator
	opcode int
}

func NewCommitList(comm netio.Communicator, opcode int) *Commitlist {
	return &Commitlist{comm, opcode}
}

func (c *Commitlist) GetCommitList() ([]string, error) {
	if err := c.comm.SendVarInt(int64(c.opcode)); err != nil {
		return nil, fmt.Errorf("cannot send opcode:\n%w", err)
	}

	err := c.comm.SendVarInt(int64(c.opcode))
	if err != nil {
		return nil, fmt.Errorf("cannot send op code:\n%w", err)
	}

	commitList, err := c.comm.ReceiveStringSlice()
	if err != nil {
		return nil, fmt.Errorf("error receiving commit list:\n%w", err)
	}
	return commitList, nil
}
