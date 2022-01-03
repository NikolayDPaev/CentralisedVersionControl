package clienthandler

import (
	"fmt"
	"io"

	"github.com/NikolayDPaev/CentralisedVersionControl/server/commit"
	"github.com/NikolayDPaev/CentralisedVersionControl/server/fileIO"
	"github.com/NikolayDPaev/CentralisedVersionControl/server/netIO"
)

func getMissingBlobIds(commit *commit.Commit) ([]string, error) {
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

func receiveBlob(reader io.Reader) error {
	blobId, err := netIO.ReceiveString(reader)
	if err != nil {
		return fmt.Errorf("error receiving blobId:\n%w", err)
	}
	file, err := fileIO.NewBlob(blobId)
	if err != nil {
		return fmt.Errorf("error creating blob:\n%w", err)
	}
	defer file.Close()

	err = netIO.ReceiveFileData(reader, file)
	if err != nil {
		return fmt.Errorf("error receiving blob:\n%w", err)
	}

	return nil
}

func saveCommit(commit *commit.Commit) error {
	commitFile, err := fileIO.NewCommit(commit.Id())
	if err != nil {
		return fmt.Errorf("error creating commit file for commit %s: %w", commit.String(), err)
	}
	defer commitFile.Close()

	if err := commit.Write(commitFile); err != nil {
		return fmt.Errorf("error saving commit %s: %w", commit.String(), err)
	}
	return nil
}

func receiveCommit(reader io.Reader, writer io.Writer) error {
	commit, err := commit.ReadCommit(reader)
	if err != nil {
		return fmt.Errorf("error receiving commit: %w", err)
	}

	missingBlobIds, err := getMissingBlobIds(commit)
	if err != nil {
		return fmt.Errorf("error getting missing blobIds from commit %s: %w", commit.String(), err)
	}

	err = netIO.SendStringSlice(missingBlobIds, writer)
	if err != nil {
		return fmt.Errorf("error sending missing blobIds from commit %s: %w", commit.String(), err)
	}

	for range missingBlobIds { // gonna receive the requested number of blobs
		if err := receiveBlob(reader); err != nil {
			return err
		}
	}
	// mutex ???
	if err := saveCommit(commit); err != nil {
		return err
	}
	return nil
}
