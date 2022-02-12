package storage

import (
	"github.com/NikolayDPaev/CentralisedVersionControl/netio"
	"github.com/NikolayDPaev/CentralisedVersionControl/server/servercommit"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Storage

// Interface that defines all methods that the server storage must
// implements.
type Storage interface {
	SendBlob(blobId string, comm netio.Communicator) error
	RecvBlob(blobId string, comm netio.Communicator) error
	BlobExists(blobId string) (bool, error)
	CommitList() []string
	OpenCommit(commitId string) (*servercommit.Commit, error)
	SaveCommit(commit *servercommit.Commit) error
	CommitSize(commitId string) (int64, error)
	CommitExists(commitId string) (bool, error)
}

// Storage implementation
type FileStorage struct{}
