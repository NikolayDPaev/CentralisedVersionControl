package files

import (
	"errors"
	"fmt"
	"io"
	"os"
)

func blobPath(blobId string) string {
	return "blobs/" + blobId[:2] + "/" + blobId[2:]
}

func BlobReader(blobId string) (io.Reader, error) {
	reader, err := os.Open(blobPath(blobId))
	if err != nil {
		return nil, fmt.Errorf("cannot open blob: %w", err)
	}
	return reader, nil
}

func NewBlobWriter(blobId string) (io.Writer, error) {
	writer, err := os.Create(blobPath(blobId))
	if err != nil {
		return nil, fmt.Errorf("cannot create blob file: %w", err)
	}
	return writer, nil
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
