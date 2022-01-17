package commands

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/NikolayDPaev/CentralisedVersionControl/client/commit"
	"github.com/NikolayDPaev/CentralisedVersionControl/client/fileIO"
	"github.com/NikolayDPaev/CentralisedVersionControl/client/netIO"
)

const DOWNLOAD_COMMIT = 1
const OK = 0

var errInvalidCommitId = errors.New("invalid commit ID")

func receiveCommit(commitId string, reader io.Reader, writer io.Writer) (*commit.Commit, error) {
	if err := netIO.SendString(commitId, writer); err != nil {
		return nil, fmt.Errorf("error sending commit ID:\n%w", err)
	}

	code, err := netIO.ReceiveVarInt(reader)
	if err != nil {
		return nil, fmt.Errorf("error no such commit on server:\n%w", err)
	}
	if code != OK {
		return nil, errInvalidCommitId
	}

	fmt.Println("Receiving commit")

	commit, err := commit.ReadCommit(commitId, reader)
	if err != nil {
		return nil, fmt.Errorf("error receiving commit:\n%w", err)
	}
	return commit, nil
}

func receiveBlob(missingFilesMap map[string]string, reader io.Reader) error {
	blobId, err := netIO.ReceiveString(reader)
	if err != nil {
		return fmt.Errorf("error receiving blobId:\n%w", err)
	}
	fileName := missingFilesMap[blobId]

	tmp, err := os.CreateTemp("", "blobTmp")
	if err != nil {
		return fmt.Errorf("error creating tmpBlob:\n%w", err)
	}
	defer tmp.Close()

	if err := netIO.ReceiveFileData(reader, tmp); err != nil {
		return fmt.Errorf("error receiving blob:\n%w", err)
	}

	if err := fileIO.DecompressFile(fileName, tmp); err != nil {
		return fmt.Errorf("error decompressing blob:\n%w", err)
	}
	return nil
}

func receiveBlobs(missingFilesMap map[string]string, reader io.Reader) error {
	for range missingFilesMap {
		err := receiveBlob(missingFilesMap, reader)
		if err != nil {
			return err
		}
	}
	return nil
}

func DownloadCommit(commitId string, reader io.Reader, writer io.Writer) (string, error) {
	fmt.Printf("Requesting commit %s\n", commitId)
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
		return "", fmt.Errorf("error getting missing files from commit:\n%w", err)
	}

	missingBlobIds := make([]string, 0, len(missingFilesMap))
	for k, _ := range missingFilesMap {
		missingBlobIds = append(missingBlobIds, k)
	}

	fmt.Printf("Requesting missing files\n")
	if err := netIO.SendStringSlice(missingBlobIds, writer); err != nil {
		return "", fmt.Errorf("error sending missing blobIds:\n%w", err)
	}

	fmt.Printf("Requesting missing objects: %d\n", len(missingFilesMap))
	if err := receiveBlobs(missingFilesMap, reader); err != nil {
		return "", fmt.Errorf("error receiving missing blobs:\n%w", err)
	}

	return "Successfuly downloaded commit", nil
}
