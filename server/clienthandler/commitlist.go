package clienthandler

import (
	"fmt"

	"github.com/NikolayDPaev/CentralisedVersionControl/server/fileIO"
	"github.com/NikolayDPaev/CentralisedVersionControl/server/netIO"
)

type CommitList struct {
	comm netIO.Communicator
}

func (c *CommitList) sendCommitList() error {
	metadataList := fileIO.CommitList()
	err := c.comm.SendVarInt(int64(len(metadataList)))
	if err != nil {
		return fmt.Errorf("could not send metadata list length:\n%w", err)
	}

	for _, entry := range metadataList {
		err := c.comm.SendString(entry)
		if err != nil {
			return fmt.Errorf("could not send metadata entry:\n%w", err)
		}
	}
	return nil
}

func (c *CommitList) Handle() error {
	return c.sendCommitList()
}
