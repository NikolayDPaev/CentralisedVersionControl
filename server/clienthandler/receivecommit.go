package clienthandler

import (
	"fmt"

	"github.com/NikolayDPaev/CentralisedVersionControl/server/commit"
	"github.com/NikolayDPaev/CentralisedVersionControl/server/fileIO"
	"github.com/NikolayDPaev/CentralisedVersionControl/server/netIO"
)

type ReceiveCommit struct {
	comm netIO.Communicator
}

func (r *ReceiveCommit) getMissingBlobIds(commit *commit.Commit) ([]string, error) {
	commitBlobIds := commit.ExtractBlobIds()
	//missingBlobIds := make([]string, len(commitBlobIds)/2)
	var missingBlobIds []string

	for _, blobId := range commitBlobIds {
		exists, err := fileIO.BlobExists(blobId)
		if err != nil {
			return nil, fmt.Errorf("cannot check existence of blob %s:\n%w", blobId, err)
		}

		if !exists {
			missingBlobIds = append(missingBlobIds, blobId)
		}
	}
	return missingBlobIds, nil
}

func (r *ReceiveCommit) receiveBlob() error {
	blobId, err := r.comm.ReceiveString()
	if err != nil {
		return fmt.Errorf("error receiving blobId:\n%w", err)
	}
	file, err := fileIO.NewBlob(blobId)
	if err != nil {
		return fmt.Errorf("error creating blob:\n%w", err)
	}
	defer file.Close()

	err = r.comm.ReceiveFileData(file)
	if err != nil {
		return err
	}

	return nil
}

func (r *ReceiveCommit) saveCommit(commit *commit.Commit) error { // !!!
	commitFile, err := fileIO.NewCommit(commit.Id())
	if err != nil {
		return fmt.Errorf("error creating commit file for commit %s: %w", commit.String(), err)
	}
	defer commitFile.Close()

	comm := netIO.NewCommunicator(100, commitFile, commitFile)
	if err := commit.Write(comm); err != nil {
		return fmt.Errorf("error saving commit %s: %w", commit.String(), err)
	}
	return nil
}

func (r *ReceiveCommit) receiveCommit() error {
	id, err := r.comm.ReceiveString()
	if err != nil {
		return fmt.Errorf("cannot read id of commit:\n%w", err)
	}

	commit, err := commit.ReadCommit(id, r.comm)
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
	if err := r.saveCommit(commit); err != nil {
		return err
	}
	return nil
}

func (r *ReceiveCommit) Handle() error {
	return r.receiveCommit()
}
