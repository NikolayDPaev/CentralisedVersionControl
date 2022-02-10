package commands

import (
	"fmt"
	"sort"

	"github.com/NikolayDPaev/CentralisedVersionControl/client/clientcommit"
	"github.com/NikolayDPaev/CentralisedVersionControl/client/fileio"
	"github.com/NikolayDPaev/CentralisedVersionControl/netio"
)

type Upload struct {
	comm     netio.Communicator
	localcpy fileio.Localcopy
	opcode   int
}

func NewUpload(comm netio.Communicator, localcpy fileio.Localcopy, opcode int) *Upload {
	return &Upload{comm, localcpy, opcode}
}

func createCommit(message, creator string, localcpy fileio.Localcopy) (*clientcommit.Commit, error) {
	paths, err := localcpy.GetPathsOfAllFiles()
	if err != nil {
		return nil, fmt.Errorf("error getting filenames for creating commit:\n%w", err)
	}

	fileSortedSlice := make([]clientcommit.CommitEntry, len(paths))
	for i, path := range paths {
		hash, err := localcpy.GetHashOfFile(path)
		if err != nil {
			return nil, err
		}
		fileSortedSlice[i] = clientcommit.CommitEntry{Hash: hash, Path: path}
	}

	sort.Slice(fileSortedSlice, func(i, j int) bool {
		return fileSortedSlice[i].Hash < fileSortedSlice[j].Hash
	})

	return &clientcommit.Commit{Message: message, Creator: creator, FileSortedSlice: fileSortedSlice}, nil
}

func (u *Upload) sendCommit(commit *clientcommit.Commit) error {
	err := u.comm.SendString(commit.Md5Hash())
	if err != nil {
		return fmt.Errorf("error sending commit id:\n%w", err)
	}

	err = u.comm.SendString(commit.Message)
	if err != nil {
		return fmt.Errorf("cannot send commit message: %w", err)
	}

	err = u.comm.SendString(commit.Creator)
	if err != nil {
		return fmt.Errorf("cannot send commit creator: %w", err)
	}

	err = u.comm.SendString(commit.GetTree())
	if err != nil {
		return fmt.Errorf("error sending commit tree: %w", err)
	}

	return nil
}

func (u *Upload) sendBlob(blobId string, commit *clientcommit.Commit) error {
	if err := u.comm.SendString(blobId); err != nil {
		return fmt.Errorf("error sending blobId %s:\n%w", blobId, err)
	}

	path, err := commit.GetBlobPath(blobId)
	if err != nil {
		return err
	}

	if err := u.localcpy.SendBlob(path, u.comm); err != nil {
		return fmt.Errorf("error sending file %s:\n%w", path, err)
	}
	return nil
}

func (u *Upload) UploadCommit(message, username string) error {
	if err := u.comm.SendVarInt(int64(u.opcode)); err != nil {
		return fmt.Errorf("cannot send opcode:\n%w", err)
	}

	commit, err := createCommit(message, username, u.localcpy)
	if err != nil {
		return fmt.Errorf("error creating commit:\n%w", err)
	}

	fmt.Printf("Sending commit %s\n", commit.Md5Hash())
	err = u.sendCommit(commit)
	if err != nil {
		return err
	}

	fmt.Println("Receiving missing blob IDs from server")
	missingBlobIds, err := u.comm.RecvStringSlice()
	if err != nil {
		return fmt.Errorf("error receiving missing blobIds:\n%w", err)
	}

	fmt.Printf("Sending %d objects\n", len(missingBlobIds))
	for _, blobId := range missingBlobIds {
		err := u.sendBlob(blobId, commit)
		if err != nil {
			return fmt.Errorf("error sending blob with id %s:\n%w", blobId, err)
		}
	}
	return nil
}
