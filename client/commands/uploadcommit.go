package commands

import (
	"fmt"
	"sort"

	"github.com/NikolayDPaev/CentralisedVersionControl/client/clientcommit"
	"github.com/NikolayDPaev/CentralisedVersionControl/client/fileio"
	"github.com/NikolayDPaev/CentralisedVersionControl/netio"
)

// Implements the Upload commit operation
type Upload struct {
	comm     netio.Communicator
	localcpy fileio.Localcopy
	opcode   int
}

func NewUpload(comm netio.Communicator, localcpy fileio.Localcopy, opcode int) *Upload {
	return &Upload{comm, localcpy, opcode}
}

// Creates a commit structure representing the current files in the working directory.
// Turns to the localcopy interface to get all paths and the hashes.
// If any of the operations fails - returns error.
func createCommit(message, creator string, localcpy fileio.Localcopy) (*clientcommit.Commit, error) {
	paths, err := localcpy.GetPathsOfAllFiles()
	if err != nil {
		return nil, fmt.Errorf("error getting filenames for creating commit: %w", err)
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

// Sends the commit struct to the server via the communicator interface.
// The commit is broken down to id, message, creator and tree string representing the blobIds and paths
// Returns error if any of the communication operations fails.
func (u *Upload) sendCommit(commit *clientcommit.Commit) error {
	err := u.comm.SendString(commit.Md5Hash())
	if err != nil {
		return fmt.Errorf("error sending commit id: %w", err)
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

// Sends blob to the server via the communicator interface.
// More specificaly sends the blobId and then turns to the localcopy interface
// to send the file contains.
func (u *Upload) sendBlob(blobId string, commit *clientcommit.Commit) error {
	if err := u.comm.SendString(blobId); err != nil {
		return fmt.Errorf("error sending blobId %s: %w", blobId, err)
	}

	path, err := commit.GetBlobPath(blobId)
	if err != nil {
		return err
	}

	if err := u.localcpy.SendBlob(path, u.comm); err != nil {
		return fmt.Errorf("error sending file %s: %w", path, err)
	}
	return nil
}

// Method that encapsulates all the logic behind upload commit operation.
// Sends opcode, then creates and sends commit and then sends the requested
// blobs that the server is missing.
// Returns error if any operation fails.
func (u *Upload) UploadCommit(message, username string) error {
	if err := u.comm.SendVarInt(int64(u.opcode)); err != nil {
		return fmt.Errorf("cannot send opcode: %w", err)
	}

	commit, err := createCommit(message, username, u.localcpy)
	if err != nil {
		return fmt.Errorf("error creating commit: %w", err)
	}

	fmt.Printf("Sending commit %s ", commit.Md5Hash())
	err = u.sendCommit(commit)
	if err != nil {
		return err
	}

	missingBlobIds, err := u.comm.RecvStringSlice()
	if err != nil {
		return fmt.Errorf("error receiving missing blobIds: %w", err)
	}

	fmt.Printf("Sending %d objects\n", len(missingBlobIds))
	for _, blobId := range missingBlobIds {
		err := u.sendBlob(blobId, commit)
		if err != nil {
			return fmt.Errorf("error sending blob with id %s: %w", blobId, err)
		}
	}
	return nil
}
