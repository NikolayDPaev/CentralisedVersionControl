package fileIO

import (
	"errors"
	"fmt"
	"os"
)

func blobPath(blobId string) string {
	return "blobs/" + blobId[:2] + "/" + blobId[2:]
}

func OpenBlob(blobId string) (*os.File, error) {
	file, err := os.Open(blobPath(blobId))
	if err != nil {
		return nil, fmt.Errorf("cannot open blob %s: %w", blobId, err)
	}
	return file, nil
}

func NewBlob(blobId string) (*os.File, error) {
	file, err := os.Create(blobPath(blobId))
	if err != nil {
		return nil, fmt.Errorf("cannot create blob file %s: %w", blobId, err)
	}
	return file, nil
}

func CheckIfExist(blobId string) (bool, error) {
	if _, err := os.Stat(blobPath(blobId)); err == nil {
		return true, nil

	} else if errors.Is(err, os.ErrNotExist) {
		return false, nil

	} else {
		return false, err
	}
}

func BlobSize(blobId string) (int64, error) {
	fileInfo, err := os.Stat(blobPath(blobId))
	if err != nil {
		return 0, fmt.Errorf("cannot get blob %s file info: %w", blobId, err)
	}

	return fileInfo.Size(), nil
}
