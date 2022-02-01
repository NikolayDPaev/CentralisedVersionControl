package clienthandler_test

import (
	"io"

	"github.com/NikolayDPaev/CentralisedVersionControl/netio"
	"github.com/NikolayDPaev/CentralisedVersionControl/server/servercommit"
	"github.com/NikolayDPaev/CentralisedVersionControl/server/storage"
)

type communicatorDefMock struct{}

func (c *communicatorDefMock) SendVarInt(num int64) error {
	return nil
}
func (c *communicatorDefMock) RecvVarInt() (int64, error) {
	return 0, nil
}
func (c *communicatorDefMock) SendString(str string) error {
	return nil
}
func (c *communicatorDefMock) RecvString() (string, error) {
	return "", nil
}
func (c *communicatorDefMock) SendStringSlice(slice []string) error {
	return nil
}
func (c *communicatorDefMock) RecvStringSlice() ([]string, error) {
	return nil, nil
}
func (c *communicatorDefMock) SendFileData(fileReader io.Reader, fileLength int64) error {
	return nil
}
func (c *communicatorDefMock) RecvFileData(fileWriter io.Writer) error {
	return nil
}

type storageDefMock struct{}

func (s *storageDefMock) OpenBlob(blobId string) (storage.StorageEntry, error) {
	return nil, nil
}
func (s *storageDefMock) SaveBlob(blobId string, comm netio.Communicator) error {
	return nil
}
func (s *storageDefMock) BlobExists(blobId string) (bool, error) {
	return false, nil
}
func (s *storageDefMock) BlobSize(blobId string) (int64, error) {
	return 0, nil
}
func (s *storageDefMock) CommitList() []string {
	return nil
}
func (s *storageDefMock) OpenCommit(commitId string) (storage.StorageEntry, error) {
	return nil, nil
}
func (s *storageDefMock) SaveCommit(commit *servercommit.Commit) error {
	return nil
}
func (s *storageDefMock) CommitSize(commitId string) (int64, error) {
	return 0, nil
}
func (s *storageDefMock) CommitExists(commitId string) (bool, error) {
	return false, nil
}
