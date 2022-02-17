package storage

import (
	"fmt"
	"os"

	"github.com/NikolayDPaev/CentralisedVersionControl/netio"
)

// blobPath returns a path from blobId
func blobPath(blobId string) (string, error) {
	if len(blobId) < 2 {
		return "", fmt.Errorf("invalid length of blobId: %s", blobId)
	}
	return "blobs/" + blobId[:2] + "/" + blobId[2:], nil
}

// blobSize returns the size on disk of a the blob with the specified blob id
func (s *FileStorage) blobSize(blobId string) (int64, error) {
	path, err := blobPath(blobId)
	if err != nil {
		return 0, err
	}
	fileInfo, err := os.Stat(path)
	if err != nil {
		return 0, fmt.Errorf("cannot get blob %s file info: %w", blobId, err)
	}

	return fileInfo.Size(), nil
}

// SendBlob sends blob with the specified blobId to the client.
// Returns error if open file or send fail data fails.
func (s *FileStorage) SendBlob(blobId string, comm netio.Communicator) error {
	path, err := blobPath(blobId)
	if err != nil {
		return err
	}
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("cannot open blob %s: %w", blobId, err)
	}
	defer file.Close()

	size, err := s.blobSize(blobId)
	if err != nil {
		return fmt.Errorf("error getting blob %s size: %w", blobId, err)
	}

	err = comm.SendFileData(file, size)
	if err != nil {
		return fmt.Errorf("error sending blob %s: %w", blobId, err)
	}

	return nil
}

// RecvBlob reads blob from the client and saves it on the disk.
func (s *FileStorage) RecvBlob(blobId string, comm netio.Communicator) error {
	path, err := blobPath(blobId)
	if err != nil {
		return err
	}

	if err := os.MkdirAll("blobs/"+blobId[:2], 0777); err != nil {
		return fmt.Errorf("cannot create blob folder %s: %w", blobId, err)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("cannot create blob file %s: %w", blobId, err)
	}
	defer file.Close()

	err = comm.RecvFileData(file)
	if err != nil {
		return err
	}
	return nil
}

// BlobExists predicate that checks if there is blob with this id on the disk.
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
