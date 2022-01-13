package clienthandler

import (
	"fmt"

	"github.com/NikolayDPaev/CentralisedVersionControl/server/commit"
	"github.com/NikolayDPaev/CentralisedVersionControl/server/fileIO"
	"github.com/NikolayDPaev/CentralisedVersionControl/server/netIO"
)

const (
	OK    = 0
	ERROR = 1
)

type SendCommit struct {
	comm    netIO.Communicator
	storage fileIO.Storage
}

func NewSendCommit(comm netIO.Communicator, storage fileIO.Storage) *ReceiveCommit {
	return &ReceiveCommit{comm, storage}
}

func (s *SendCommit) sendCommitData(commitId string) error {
	commitFile, err := s.storage.OpenCommit(commitId)
	if err != nil {
		return fmt.Errorf("error opening commit file of commit %s: %s", commitId, err)
	}
	defer commitFile.Close()

	commit, err := commit.ReadCommit(commitId, s.comm)
	if err != nil {
		return fmt.Errorf("error reading commit file %s: %s", commitId, err)
	}

	err = commit.Write(s.comm)
	if err != nil {
		return fmt.Errorf("error sending commit %s: %s", commitId, err)
	}
	return nil
}

func (s *SendCommit) sendBlob(blobId string) error {
	file, err := s.storage.OpenBlob(blobId)
	if err != nil {
		return fmt.Errorf("error opening blob %s:\n%w", blobId, err)
	}
	defer file.Close()

	if err := s.comm.SendString(blobId); err != nil {
		return fmt.Errorf("error sending blobId %s:\n%w", blobId, err)
	}

	size, err := s.storage.BlobSize(blobId)
	if err != nil {
		return fmt.Errorf("error getting blob %s size:\n%w", blobId, err)
	}

	err = s.comm.SendFileData(file, size)
	if err != nil {
		return fmt.Errorf("error sending blob %s:\n%w", blobId, err)
	}

	return nil
}

func (s *SendCommit) validateCommitId(commitId string) (bool, error) {
	exists, err := s.storage.CommitExists(commitId)
	if err != nil {
		return false, err
	}
	if exists {
		err = s.comm.SendVarInt(OK)
		if err != nil {
			return false, fmt.Errorf("error sending validating commit message:\n%w", err)
		}
		return true, nil
	}
	err = s.comm.SendVarInt(ERROR)
	if err != nil {
		return false, fmt.Errorf("error sending validating commit message:\n%w", err)
	}

	return false, nil
}

func (s *SendCommit) sendCommit() error {
	commitId, err := s.comm.ReceiveString()
	if err != nil {
		return fmt.Errorf("error reading commit id:\n%w", err)
	}

	validId, err := s.validateCommitId(commitId)
	if err != nil {
		return fmt.Errorf("error validating commit id:\n%w", err)
	}
	if !validId {
		return nil
	}

	err = s.sendCommitData(commitId)
	if err != nil {
		return err
	}

	blobIdsForSend, err := s.comm.ReceiveStringSlice()
	if err != nil {
		return fmt.Errorf("error getting blob ids for send:\n%w", err)
	}

	for _, blobId := range blobIdsForSend { // send the requested number of blobs
		err = s.sendBlob(blobId)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *SendCommit) Handle() error {
	return s.sendCommit()
}
