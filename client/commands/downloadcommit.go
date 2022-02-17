package commands

import (
	"errors"
	"fmt"

	"github.com/NikolayDPaev/CentralisedVersionControl/client/clientcommit"
	"github.com/NikolayDPaev/CentralisedVersionControl/client/fileio"
	"github.com/NikolayDPaev/CentralisedVersionControl/netio"
)

var ErrInvalidCommitId = errors.New("invalid commit ID")

// Download implements the download commit operation
type Download struct {
	comm     netio.Communicator
	localcpy fileio.Localcopy
	opcode   int
	okcode   int
}

func NewDownload(comm netio.Communicator, localcpy fileio.Localcopy, opcode, okcode int) *Download {
	return &Download{comm, localcpy, opcode, okcode}
}

// readCommit deserializes a commit from the communicator.
// Returns new clientcommit or error.
func (d *Download) readCommit(id string, comm netio.Communicator) (*clientcommit.Commit, error) {
	receivedId, err := comm.RecvString()
	if err != nil || id != receivedId {
		return nil, fmt.Errorf("error receiving id of commit: %w", err)
	}

	message, err := comm.RecvString()
	if err != nil {
		return nil, fmt.Errorf("cannot read message string of commit: %w", err)
	}

	creator, err := comm.RecvString()
	if err != nil {
		return nil, fmt.Errorf("cannot read creator string of commit: %w", err)
	}

	strTree, err := comm.RecvString()
	if err != nil {
		return nil, fmt.Errorf("cannot read tree string of commit: %w", err)
	}

	fileSortedSlice, err := clientcommit.GetSortedSlice(strTree)
	if err != nil {
		return nil, err
	}
	commit := &clientcommit.Commit{Message: message, Creator: creator, FileSortedSlice: fileSortedSlice}

	commitHash := commit.Md5Hash()
	if id != commitHash {
		return nil, fmt.Errorf("mismatched hash values: expected: %s, actual: %s", id, commitHash)
	}
	return commit, nil
}

// receiveCommit requests specific commit from the server.
// Returns no such commit error if the server responds with error code.
//
// On success receives commit data that represents the commit, but not the blobs.
func (d *Download) receiveCommit(commitId string) (*clientcommit.Commit, error) {
	if err := d.comm.SendString(commitId); err != nil {
		return nil, fmt.Errorf("error sending commit ID: %w", err)
	}

	code, err := d.comm.RecvVarInt()
	if err != nil {
		return nil, fmt.Errorf("error no such commit on server: %w", err)
	}
	if code != int64(d.okcode) {
		return nil, ErrInvalidCommitId
	}

	commit, err := d.readCommit(commitId, d.comm)
	if err != nil {
		return nil, fmt.Errorf("error receiving commit: %w", err)
	}
	return commit, nil
}

// getMissingFiles returns map[hash]path that contains the blobs
// that are part of the commit but are missing locally.
func (d *Download) getMissingFiles(c *clientcommit.Commit) (map[string]string, error) {
	missingFileMap := make(map[string]string, len(c.FileSortedSlice)/2)

	for _, entry := range c.FileSortedSlice {
		exists, err := d.localcpy.FileWithHashExists(entry.Path, entry.Hash)
		if err != nil {
			return nil, err
		}

		if !exists {
			missingFileMap[entry.Hash] = entry.Path
		}
	}

	return missingFileMap, nil
}

// receiveBlobs invokes ReceiveBlob method of the localcopy package on every entry in the missingFilesMap.
func (d *Download) receiveBlobs(missingFilesMap map[string]string) error {
	for range missingFilesMap {
		blobId, err := d.comm.RecvString()
		if err != nil {
			return fmt.Errorf("error receiving blobId: %w", err)
		}
		fileName := missingFilesMap[blobId]
		d.localcpy.ReceiveBlob(fileName, d.comm)
	}
	return nil
}

// DownloadCommit is the package main method that encapsulates all the logic behind downlonading a commit.
// Requests the commit, reads it, requests the locally missing blobs and then downloads them.
// If any operation fails - returns error
func (d *Download) DownloadCommit(commitId string) error {
	if err := d.comm.SendVarInt(int64(d.opcode)); err != nil {
		return fmt.Errorf("cannot send opcode: %w", err)
	}

	fmt.Printf("Requesting commit %s\n", commitId)
	commit, err := d.receiveCommit(commitId)
	if err != nil {
		return err
	}

	if err := d.localcpy.CleanOtherFiles(commit.GetSetOfPaths()); err != nil {
		return fmt.Errorf("error cleaning other files: %w", err)
	}

	missingFilesMap, err := d.getMissingFiles(commit)
	if err != nil {
		return fmt.Errorf("error getting missing files from commit: %w", err)
	}

	var missingBlobIds []string
	for k := range missingFilesMap {
		missingBlobIds = append(missingBlobIds, k)
	}

	if err := d.comm.SendStringSlice(missingBlobIds); err != nil {
		return fmt.Errorf("error sending missing blobIds: %w", err)
	}

	fmt.Printf("Receiving missing objects: %d\n", len(missingFilesMap))
	if err := d.receiveBlobs(missingFilesMap); err != nil {
		return fmt.Errorf("error receiving missing blobs: %w", err)
	}

	fmt.Println("Commit downloaded successfuly!")

	return nil
}
