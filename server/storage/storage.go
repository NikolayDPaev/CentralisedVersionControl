package storage

import (
	"github.com/NikolayDPaev/CentralisedVersionControl/netio"
	"github.com/NikolayDPaev/CentralisedVersionControl/server/servercommit"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Storage

// Interface that defines all methods that the server storage must
// implements.
type Storage interface {
	// Sends blob with the specified blobId via the communicator.
	// First sends the size, then the data.
	SendBlob(blobId string, comm netio.Communicator) error

	// Reads blob with the specified blobId via the communicator.
	// First reads the size, then the data.
	RecvBlob(blobId string, comm netio.Communicator) error

	// Checks for existence of blob with the specified id.
	BlobExists(blobId string) (bool, error)

	// Returns a slice with string representation of the meta data
	// of all commits in the storage.
	CommitList() []string

	// Returns a Commit struct that represents the commit with the provided id.
	OpenCommit(commitId string) (*servercommit.Commit, error)

	// Writes the commit data from the Commit struct to the storage.
	SaveCommit(commit *servercommit.Commit) error

	// Checks if commmit with the provided id is present in the storage.
	CommitExists(commitId string) (bool, error)
}

// Storage implementation that relies on the file system to store
// the commits and the blobs as files.
type FileStorage struct{}
