package commands

import (
	"errors"
	"fmt"

	"github.com/NikolayDPaev/CentralisedVersionControl/client/clientcommit"
	"github.com/NikolayDPaev/CentralisedVersionControl/client/fileio"
	"github.com/NikolayDPaev/CentralisedVersionControl/netio"
)

var ErrInvalidCommitId = errors.New("invalid commit ID")

type Download struct {
	comm     netio.Communicator
	localcpy fileio.Localcopy
	opcode   int
	okcode   int
}

func NewDownload(comm netio.Communicator, localcpy fileio.Localcopy, opcode, okcode int) *Download {
	return &Download{comm, localcpy, opcode, okcode}
}

func ReadCommit(id string, comm netio.Communicator) (*clientcommit.Commit, error) {
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

	fmt.Println("Receiving commit")

	commit, err := ReadCommit(commitId, d.comm)
	if err != nil {
		return nil, fmt.Errorf("error receiving commit: %w", err)
	}
	return commit, nil
}

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

	missingFilesMap, err := commit.GetMissingFiles(d.localcpy)
	if err != nil {
		return fmt.Errorf("error getting missing files from commit: %w", err)
	}

	var missingBlobIds []string
	for k := range missingFilesMap {
		missingBlobIds = append(missingBlobIds, k)
	}

	fmt.Printf("Requesting missing files\n")
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
