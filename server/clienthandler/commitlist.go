package clienthandler

import (
	"fmt"

	"github.com/NikolayDPaev/CentralisedVersionControl/server/fileIO"
	"github.com/NikolayDPaev/CentralisedVersionControl/server/netIO"
)

type CommitList struct {
	comm    netIO.Communicator
	storage fileIO.Storage
}

func NewCommitList(comm netIO.Communicator, storage fileIO.Storage) *CommitList {
	return &CommitList{comm, storage}
}

func (c *CommitList) sendCommitList() error {
	metadataList := c.storage.CommitList()
	err := c.comm.SendStringSlice(metadataList)
	if err != nil {
		return fmt.Errorf("could not send metadata list:\n%w", err)
	}
	return nil
}

func (c *CommitList) Handle() error {
	return c.sendCommitList()
}
