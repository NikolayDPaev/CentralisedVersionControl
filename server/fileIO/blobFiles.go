package fileIO

import (
	"fmt"
	"os"
)

func blobPath(blobId string) (string, error) {
	if len(blobId) < 2 {
		return "", fmt.Errorf("invalid length of blobId:\n%s", blobId)
	}
	return "blobs/" + blobId[:2] + "/" + blobId[2:], nil
}

func OpenBlob(blobId string) (*os.File, error) {
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

func NewBlob(blobId string) (*os.File, error) {
	path, err := blobPath(blobId)
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll("blobs/"+blobId[:2], 0777); err != nil {
		return nil, fmt.Errorf("cannot create blob folder %s:\n%w", blobId, err)
	}

	file, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("cannot create blob file %s:\n%w", blobId, err)
	}
	return file, nil
}

func BlobExists(blobId string) (bool, error) {
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

func BlobSize(blobId string) (int64, error) {
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
