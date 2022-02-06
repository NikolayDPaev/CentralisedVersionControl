package clienthandler

import (
	"fmt"

	"github.com/NikolayDPaev/CentralisedVersionControl/netio"
	"github.com/NikolayDPaev/CentralisedVersionControl/server/servercommit"
	"github.com/NikolayDPaev/CentralisedVersionControl/server/storage"
)

type ReceiveCommit struct {
	comm    netio.Communicator
	storage storage.Storage
}

func NewReceiveCommit(comm netio.Communicator, storage storage.Storage) *ReceiveCommit {
	return &ReceiveCommit{comm, storage}
}

func (r *ReceiveCommit) getMissingBlobIds(commit *servercommit.Commit) ([]string, error) {
	commitBlobIds := commit.ExtractBlobIds()
	//missingBlobIds := make([]string, len(commitBlobIds)/2)
	var missingBlobIds []string

	for _, blobId := range commitBlobIds {
		exists, err := r.storage.BlobExists(blobId)
		if err != nil {
			return nil, fmt.Errorf("cannot check existence of blob %s: %w", blobId, err)
		}

		if !exists {
			missingBlobIds = append(missingBlobIds, blobId)
		}
	}
	return missingBlobIds, nil
}

func (r *ReceiveCommit) receiveBlob() error {
	blobId, err := r.comm.RecvString()
	if err != nil {
		return fmt.Errorf("error receiving blobId: %w", err)
	}
	err = r.storage.SaveBlob(blobId, r.comm)
	if err != nil {
		return fmt.Errorf("error creating blob: %w", err)
	}

	return nil
}

func (r *ReceiveCommit) receiveCommit() error {
	id, err := r.comm.RecvString()
	if err != nil {
		return fmt.Errorf("cannot read id of commit: %w", err)
	}

	commit, err := servercommit.NewCommitFrom(id, r.comm)
	if err != nil {
		return fmt.Errorf("error receiving commit: %w", err)
	}

	missingBlobIds, err := r.getMissingBlobIds(commit)
	if err != nil {
		return fmt.Errorf("error getting missing blobIds from commit %s: %w", commit.String(), err)
	}

	err = r.comm.SendStringSlice(missingBlobIds)
	if err != nil {
		return fmt.Errorf("error sending missing blobIds from commit %s: %w", commit.String(), err)
	}

	for range missingBlobIds { // gonna receive the requested number of blobs
		if err := r.receiveBlob(); err != nil {
			return err
		}
	}
	// mutex ???
	if err := r.storage.SaveCommit(commit); err != nil {
		return err
	}
	return nil
}

func (r *ReceiveCommit) Handle() error {
	return r.receiveCommit()
}
