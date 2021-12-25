package commands

import (
	"errors"
	"fmt"
	"io"

	"github.com/NikolayDPaev/CentralisedVersionControl/client/commit"
	"github.com/NikolayDPaev/CentralisedVersionControl/client/netIO"
)

const (
	GET_COMMIT_LIST = 0
	DOWNLOAD_COMMIT = 1
	UPLOAD_COMMIT   = 2
	EMPTY_REQUEST   = 3
	OK              = 0
)

var errInvalidCommitId = errors.New("invalid commit ID")

func GetCommitList(reader io.Reader, writer io.Writer) ([]string, error) {
	err := netIO.SendVarInt(GET_COMMIT_LIST, writer)
	if err != nil {
		return nil, fmt.Errorf("cannot send op code: %w", err)
	}

	commitList, err := netIO.ReceiveStringSlice(reader)
	if err != nil {
		return nil, fmt.Errorf("error receiving commit list: %w", err)
	}
	return commitList, nil
}

func receiveCommit(commitId string, reader io.Reader, writer io.Writer) (*commit.Commit, error) {
	if err := netIO.SendString(commitId, writer); err != nil {
		return nil, fmt.Errorf("error sending commit ID: %w", err)
	}

	code, err := netIO.ReceiveVarInt(reader)
	if err != nil {
		return nil, fmt.Errorf("error receiving success code: %w", err)
	}
	if code != OK {
		return nil, errInvalidCommitId
	}

	commit, err := commit.ReadCommit(reader)
	if err != nil {
		return nil, fmt.Errorf("error receiving commit: %w", err)
	}
	return commit, nil
}

func receiveBlobs(missingFilesMap map[string]string, reader io.Reader) error {

	return nil
}

func DownloadCommit(commitId string, reader io.Reader, writer io.Writer) (string, error) {
	commit, err := receiveCommit(commitId, reader, writer)
	if err != nil {
		if errors.Is(err, errInvalidCommitId) {
			return "Invalid commit Id", nil
		} else {
			return "", nil
		}
	}

	// delete every other file

	missingFilesMap, err := commit.GetMissingFiles()
	if err != nil {
		return "", fmt.Errorf("error getting missing files from commit: %w", err)
	}

	missingBlobIds := make([]string, 0, len(missingFilesMap))
	for k := range missingFilesMap {
		missingBlobIds = append(missingBlobIds, k)
	}

	if err := netIO.SendStringSlice(missingBlobIds, writer); err != nil {
		return "", fmt.Errorf("error sending missing blobIds: %w", err)
	}

	if err := receiveBlobs(missingFilesMap, reader); err != nil {
		return "", fmt.Errorf("error receiving missing blobIds: %w", err)
	}

	return "Successfuly downloaded commit", nil
}
