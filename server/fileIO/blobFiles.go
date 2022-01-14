package fileIO

import (
	"fmt"
	"os"

	"github.com/NikolayDPaev/CentralisedVersionControl/server/netIO"
)

func blobPath(blobId string) (string, error) {
	if len(blobId) < 2 {
		return "", fmt.Errorf("invalid length of blobId: %s", blobId)
	}
	return "blobs/" + blobId[:2] + "/" + blobId[2:], nil
}

func (s *FileStorage) OpenBlob(blobId string) (StorageEntry, error) {
	path, err := blobPath(blobId)
	if err != nil {
		return nil, err
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("cannot open blob %s:\n%w", blobId, err)
	}
	return file, nil
}

func (s *FileStorage) SaveBlob(blobId string, comm netIO.Communicator) error {
	path, err := blobPath(blobId)
	if err != nil {
		return err
	}

	if err := os.MkdirAll("blobs/"+blobId[:2], 0777); err != nil {
		return fmt.Errorf("cannot create blob folder %s:\n%w", blobId, err)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("cannot create blob file %s:\n%w", blobId, err)
	}
	defer file.Close()

	err = comm.ReceiveFileData(file)
	if err != nil {
		return err
	}
	return nil
}

func (s *FileStorage) BlobExists(blobId string) (bool, error) {
	path, err := blobPath(blobId)
	if err != nil {
		return false, err
	}
	b, err := fileExists(path)
	if err != nil {
		return false, err
	}
	return b, nil
}

func (s *FileStorage) BlobSize(blobId string) (int64, error) {
	path, err := blobPath(blobId)
	if err != nil {
		return 0, err
	}
	fileInfo, err := os.Stat(path)
	if err != nil {
		return 0, fmt.Errorf("cannot get blob %s file info:\n%w", blobId, err)
	}

	return fileInfo.Size(), nil
}
