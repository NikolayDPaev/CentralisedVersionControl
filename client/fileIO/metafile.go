package fileIO

import (
	"fmt"
	"os"
)

const (
	METAFILE_NAME = ".cvc"
)

func NewMetaFile() (*os.File, error) {
	file, err := os.Create(METAFILE_NAME)
	if err != nil {
		return nil, fmt.Errorf("cannot create metafile: %w", err)
	}
	return file, nil
}

func OpenMetaFile() (*os.File, error) {
	file, err := os.Open(METAFILE_NAME)
	if err != nil {
		return nil, fmt.Errorf("cannot open metafile: %w", err)
	}
	return file, nil
}
