package storage

import (
	"github.com/NikolayDPaev/CentralisedVersionControl/netio"
	"github.com/NikolayDPaev/CentralisedVersionControl/server/servercommit"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Storage

// Storage is an interface that defines all methods that the server storage must
// implements.
type Storage interface {
	// SendBlob sends blob with the specified blobId via the communicator.
	// First sends the size, then the data.
	SendBlob(blobId string, comm netio.Communicator) error

	// RecvBlob reads blob with the specified blobId via the communicator.
	// First reads the size, then the data.
	RecvBlob(blobId string, comm netio.Communicator) error

	// BlobExists checks for existence of blob with the specified id.
	BlobExists(blobId string) (bool, error)

	// CommitList returns a slice with string representation of the meta data
	// of all commits in the storage.
	CommitList() []string

	// OpenCommit returns a Commit struct that represents the commit with the provided id.
	OpenCommit(commitId string) (*servercommit.Commit, error)

	// SaveCommit writes the commit data from the Commit struct to the storage.
	SaveCommit(commit *servercommit.Commit) error

	// CommitExists checks if commmit with the provided id is present in the storage.
	CommitExists(commitId string) (bool, error)
}

// FileStorage is implementation of Storage that relies on the file system to store
// the commits and the blobs as files.
type FileStorage struct{}
