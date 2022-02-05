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

func (d *Download) receiveCommit(commitId string) (*clientcommit.Commit, error) {
	if err := d.comm.SendString(commitId); err != nil {
		return nil, fmt.Errorf("error sending commit ID:\n%w", err)
	}

	code, err := d.comm.RecvVarInt()
	if err != nil {
		return nil, fmt.Errorf("error no such commit on server:\n%w", err)
	}
	if code != int64(d.okcode) {
		return nil, ErrInvalidCommitId
	}

	fmt.Println("Receiving commit")

	commit, err := clientcommit.ReadCommit(commitId, d.comm)
	if err != nil {
		return nil, fmt.Errorf("error receiving commit:\n%w", err)
	}
	return commit, nil
}

func (d *Download) receiveBlob(missingFilesMap map[string]string) error {
	blobId, err := d.comm.RecvString()
	if err != nil {
		return fmt.Errorf("error receiving blobId:\n%w", err)
	}
	fileName := missingFilesMap[blobId]

	return d.localcpy.ReceiveBlob(fileName, d.comm)
}

func (d *Download) receiveBlobs(missingFilesMap map[string]string) error {
	for range missingFilesMap {
		err := d.receiveBlob(missingFilesMap)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *Download) DownloadCommit(commitId string) error {
	if err := d.comm.SendVarInt(int64(d.opcode)); err != nil {
		return fmt.Errorf("cannot send opcode:\n%w", err)
	}

	fmt.Printf("Requesting commit %s\n", commitId)
	commit, err := d.receiveCommit(commitId)
	if err != nil {
		return err
	}

	if err := d.localcpy.CleanOtherFiles(commit.GetSetOfPaths()); err != nil {
		return fmt.Errorf("error cleaning other files:\n%w", err)
	}

	missingFilesMap, err := commit.GetMissingFiles(d.localcpy)
	if err != nil {
		return fmt.Errorf("error getting missing files from commit:\n%w", err)
	}

	var missingBlobIds []string
	for k := range missingFilesMap {
		missingBlobIds = append(missingBlobIds, k)
	}

	fmt.Printf("Requesting missing files\n")
	if err := d.comm.SendStringSlice(missingBlobIds); err != nil {
		return fmt.Errorf("error sending missing blobIds:\n%w", err)
	}

	fmt.Printf("Requesting missing objects: %d\n", len(missingFilesMap))
	if err := d.receiveBlobs(missingFilesMap); err != nil {
		return fmt.Errorf("error receiving missing blobs:\n%w", err)
	}

	return nil
}
