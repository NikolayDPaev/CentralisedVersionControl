package clienthandler

import (
	"fmt"

	"github.com/NikolayDPaev/CentralisedVersionControl/netio"
	"github.com/NikolayDPaev/CentralisedVersionControl/server/storage"
)

// CommitList implements the commit list operation.
type CommitList struct {
	comm    netio.Communicator
	storage storage.Storage
}

func NewCommitList(comm netio.Communicator, storage storage.Storage) *CommitList {
	return &CommitList{comm, storage}
}

// sendCommitList sends commitList to the client as a string slice.
func (c *CommitList) sendCommitList() error {
	metadataList := c.storage.CommitList()
	err := c.comm.SendStringSlice(metadataList)
	if err != nil {
		return fmt.Errorf("could not send metadata list: %w", err)
	}
	return nil
}

func (c *CommitList) Handle() error {
	return c.sendCommitList()
}
