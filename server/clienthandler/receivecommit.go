package clienthandler

import (
	"fmt"

	"github.com/NikolayDPaev/CentralisedVersionControl/netio"
	"github.com/NikolayDPaev/CentralisedVersionControl/server/servercommit"
	"github.com/NikolayDPaev/CentralisedVersionControl/server/storage"
)

// ReceiveCommit implements the receive commit operation.
type ReceiveCommit struct {
	comm    netio.Communicator
	storage storage.Storage
}

func NewReceiveCommit(comm netio.Communicator, storage storage.Storage) *ReceiveCommit {
	return &ReceiveCommit{comm, storage}
}

// getMissingBlobIds returns slice with the blob ids from the commit that are missing on the server and
// must be requested or error if BLobExists fails
func (r *ReceiveCommit) getMissingBlobIds(commit *servercommit.Commit) ([]string, error) {
	commitBlobIds := commit.ExtractBlobIds()
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

// receiveBlob receives a blob from the client and saves it on the storage
func (r *ReceiveCommit) receiveBlob() error {
	blobId, err := r.comm.RecvString()
	if err != nil {
		return fmt.Errorf("error receiving blobId: %w", err)
	}
	err = r.storage.RecvBlob(blobId, r.comm)
	if err != nil {
		return fmt.Errorf("error creating blob: %w", err)
	}

	return nil
}

// receiveCommit provides the logic behind the receiveCommit operation.
// Reads the commit id, deserializes the commit,
// extracts the blob ids of the blobs that are missing,
// sends them to the client, receives the blobs and then
// saves the commit.
//
// Because the commit file is saved after all blobs are
// downloaded, other goroutines performing send commit operation
// will not see the current commit before it is ready, so there is
// not a possibility of race condtion.
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

	for range missingBlobIds { // receive the requested number of blobs
		if err := r.receiveBlob(); err != nil {
			return err
		}
	}

	if err := r.storage.SaveCommit(commit); err != nil {
		return err
	}
	return nil
}

func (r *ReceiveCommit) Handle() error {
	return r.receiveCommit()
}
