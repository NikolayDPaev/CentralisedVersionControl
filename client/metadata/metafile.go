package metadata

import (
	"fmt"
	"os"
)

func newMetaFile(metafileName string) (*os.File, error) {
	file, err := os.Create(metafileName)
	if err != nil {
		return nil, fmt.Errorf("cannot create metafile: %w", err)
	}
	return file, nil
}

func openMetaFile(metafileName string) (*os.File, error) {
	file, err := os.Open(metafileName)
	if err != nil {
		return nil, fmt.Errorf("cannot open metafile: %w", err)
	}
	return file, nil
}
