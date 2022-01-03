package clienthandler

import (
	"fmt"
	"io"

	"github.com/NikolayDPaev/CentralisedVersionControl/server/fileIO"
	"github.com/NikolayDPaev/CentralisedVersionControl/server/netIO"
)

const (
	OK    = 0
	ERROR = 1
)

func sendCommitData(commitId string, writer io.Writer) error {
	commitFile, err := fileIO.OpenCommit(commitId)
	if err != nil {
		return fmt.Errorf("error opening commit file of commit %s: %s", commitId, err)
	}
	defer commitFile.Close()

	commitFileSize, err := fileIO.CommitSize(commitId)
	if err != nil {
		return fmt.Errorf("error getting commit file size of commit %s: %s", commitId, err)
	}

	err = netIO.SendFileData(commitFile, commitFileSize, writer)
	if err != nil {
		return fmt.Errorf("error sending commit %s: %s", commitId, err)
	}
	return nil
}

func sendBlob(blobId string, writer io.Writer) error {
	file, err := fileIO.OpenBlob(blobId)
	if err != nil {
		return fmt.Errorf("error opening blob %s:\n%w", blobId, err)
	}
	defer file.Close()

	if err := netIO.SendString(blobId, writer); err != nil {
		return fmt.Errorf("error sending blobId %s:\n%w", blobId, err)
	}

	size, err := fileIO.BlobSize(blobId)
	if err != nil {
		return fmt.Errorf("error getting blob %s size:\n%w", blobId, err)
	}

	err = netIO.SendFileData(file, size, writer)
	if err != nil {
		return fmt.Errorf("error sending blob %s:\n%w", blobId, err)
	}

	return nil
}

func validateCommitId(commitId string, writer io.Writer) (bool, error) {
	exists, err := fileIO.CommitExists(commitId)
	if err != nil {
		return false, err
	}
	if exists {
		netIO.SendVarInt(OK, writer)
		return true, nil
	}
	netIO.SendVarInt(ERROR, writer)
	return false, nil
}

func sendCommit(reader io.Reader, writer io.Writer) error {
	commitId, err := netIO.ReceiveString(reader)
	if err != nil {
		return fmt.Errorf("error reading commit id:\n%w", err)
	}

	validId, err := validateCommitId(commitId, writer)
	if err != nil {
		return fmt.Errorf("error validating commit id:\n%w", err)
	}
	if !validId {
		return nil
	}

	err = sendCommitData(commitId, writer)
	if err != nil {
		return err
	}

	blobIdsForSend, err := netIO.ReceiveStringSlice(reader)
	if err != nil {
		return fmt.Errorf("error getting blob ids for send:\n%w", err)
	}

	for _, blobId := range blobIdsForSend { // send the requested number of blobs
		err = sendBlob(blobId, writer)
		if err != nil {
			return err
		}
	}
	return nil
}
