package storage

import (
	"github.com/NikolayDPaev/CentralisedVersionControl/netio"
	"github.com/NikolayDPaev/CentralisedVersionControl/server/servercommit"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//counterfeiter:generate . StorageEntry
type StorageEntry interface {
	Read(b []byte) (n int, err error)
	Write(b []byte) (n int, err error)
	Close() error
}

//counterfeiter:generate . Storage
type Storage interface {
	OpenBlob(blobId string) (StorageEntry, error)
	SaveBlob(blobId string, comm netio.Communicator) error
	BlobExists(blobId string) (bool, error)
	BlobSize(blobId string) (int64, error)
	CommitList() []string
	OpenCommit(commitId string) (*servercommit.Commit, error)
	SaveCommit(commit *servercommit.Commit) error
	CommitSize(commitId string) (int64, error)
	CommitExists(commitId string) (bool, error)
}

type FileStorage struct{}
