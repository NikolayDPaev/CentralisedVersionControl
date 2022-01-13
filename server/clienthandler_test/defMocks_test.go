package clienthandler_test

import (
	"io"

	"github.com/NikolayDPaev/CentralisedVersionControl/server/fileIO"
)

type communicatorDefMock struct{}

func (c *communicatorDefMock) SendVarInt(num int64) error {
	return nil
}
func (c *communicatorDefMock) ReceiveVarInt() (int64, error) {
	return 0, nil
}
func (c *communicatorDefMock) SendString(str string) error {
	return nil
}
func (c *communicatorDefMock) ReceiveString() (string, error) {
	return "", nil
}
func (c *communicatorDefMock) SendStringSlice(slice []string) error {
	return nil
}
func (c *communicatorDefMock) ReceiveStringSlice() ([]string, error) {
	return nil, nil
}
func (c *communicatorDefMock) SendFileData(fileReader io.Reader, fileLength int64) error {
	return nil
}
func (c *communicatorDefMock) ReceiveFileData(fileWriter io.Writer) error {
	return nil
}

type storageDefMock struct{}

func (s *storageDefMock) OpenBlob(blobId string) (fileIO.StorageEntry, error) {
	return nil, nil
}
func (s *storageDefMock) NewBlob(blobId string) (fileIO.StorageEntry, error) {
	return nil, nil
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
func (s *storageDefMock) OpenCommit(commitId string) (fileIO.StorageEntry, error) {
	return nil, nil
}
func (s *storageDefMock) NewCommit(commitId string) (fileIO.StorageEntry, error) {
	return nil, nil
}
func (s *storageDefMock) CommitSize(commitId string) (int64, error) {
	return 0, nil
}
func (s *storageDefMock) CommitExists(commitId string) (bool, error) {
	return false, nil
}
