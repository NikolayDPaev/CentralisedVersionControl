package fileIO

type StorageEntry interface {
	Read(b []byte) (n int, err error)
	Write(b []byte) (n int, err error)
	Close() error
}

type Storage interface {
	OpenBlob(blobId string) (StorageEntry, error)
	NewBlob(blobId string) (StorageEntry, error)
	BlobExists(blobId string) (bool, error)
	BlobSize(blobId string) (int64, error)
	CommitList() []string
	OpenCommit(commitId string) (StorageEntry, error)
	NewCommit(commitId string) (StorageEntry, error)
	CommitSize(commitId string) (int64, error)
	CommitExists(commitId string) (bool, error)
}

type FileStorage struct{}
