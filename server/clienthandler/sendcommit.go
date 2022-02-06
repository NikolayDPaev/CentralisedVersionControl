package clienthandler

import (
	"fmt"

	"github.com/NikolayDPaev/CentralisedVersionControl/netio"
	"github.com/NikolayDPaev/CentralisedVersionControl/server/storage"
)

const (
	OK    = 0
	ERROR = 1
)

type SendCommit struct {
	comm    netio.Communicator
	storage storage.Storage
}

func NewSendCommit(comm netio.Communicator, storage storage.Storage) *SendCommit {
	return &SendCommit{comm, storage}
}

func (s *SendCommit) sendCommitData(commitId string) error {
	commit, err := s.storage.OpenCommit(commitId)
	if err != nil {
		return fmt.Errorf("error opening commit file of commit %s: %s", commitId, err)
	}

	err = s.comm.SendString(commit.Id)
	if err != nil {
		return fmt.Errorf("error sending commit id %s: %s", commitId, err)
	}

	err = commit.WriteTo(s.comm)
	if err != nil {
		return fmt.Errorf("error sending commit %s: %s", commitId, err)
	}
	return nil
}

func (s *SendCommit) sendBlob(blobId string) error {
	file, err := s.storage.OpenBlob(blobId)
	if err != nil {
		return fmt.Errorf("error opening blob %s: %w", blobId, err)
	}
	defer file.Close()

	if err := s.comm.SendString(blobId); err != nil {
		return fmt.Errorf("error sending blobId %s: %w", blobId, err)
	}

	size, err := s.storage.BlobSize(blobId)
	if err != nil {
		return fmt.Errorf("error getting blob %s size: %w", blobId, err)
	}

	err = s.comm.SendFileData(file, size)
	if err != nil {
		return fmt.Errorf("error sending blob %s: %w", blobId, err)
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
	commitId, err := s.comm.RecvString()
	if err != nil {
		return fmt.Errorf("error reading commit id: %w", err)
	}

	validId, err := s.validateCommitId(commitId)
	if err != nil {
		return fmt.Errorf("error validating commit id: %w", err)
	}
	if !validId {
		return nil
	}

	err = s.sendCommitData(commitId)
	if err != nil {
		return err
	}

	blobIdsForSend, err := s.comm.RecvStringSlice()
	if err != nil {
		return fmt.Errorf("error getting blob ids for send: %w", err)
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
