package fileIO

import (
	"github.com/NikolayDPaev/CentralisedVersionControl/server/commit"
	"github.com/NikolayDPaev/CentralisedVersionControl/server/netIO"
)

type StorageEntry interface {
	Read(b []byte) (n int, err error)
	Write(b []byte) (n int, err error)
	Close() error
}

type Storage interface {
	OpenBlob(blobId string) (StorageEntry, error)
	SaveBlob(blobId string, comm netIO.Communicator) error
	BlobExists(blobId string) (bool, error)
	BlobSize(blobId string) (int64, error)
	CommitList() []string
	OpenCommit(commitId string) (StorageEntry, error)
	SaveCommit(commit *commit.Commit) error
	CommitSize(commitId string) (int64, error)
	CommitExists(commitId string) (bool, error)
}

type FileStorage struct{}
