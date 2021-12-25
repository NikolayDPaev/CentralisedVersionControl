package fileIO

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"os"
)

func fileExists(filePath string) (bool, error) {
	if _, err := os.Stat(filePath); err == nil {
		return true, nil

	} else if errors.Is(err, os.ErrNotExist) {
		return false, nil

	} else {
		return false, err
	}
}

func getHashOfFile(filepath string) (string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "", fmt.Errorf("error opening %s: %w", filepath, err)
	}
	defer file.Close()

	hash := md5.New()
	_, err = io.Copy(hash, file)

	if err != nil {
		return "", fmt.Errorf("error getting hash of %s: %w", filepath, err)
	}

	return string(hash.Sum(nil)), nil
}

func FileWithHashExists(filepath string, hash string) (bool, error) {
	fileExists, err := fileExists(filepath)
	if err != nil {
		return false, fmt.Errorf("error checking if file exists: %w", err)
	}
	if !fileExists {
		return false, nil
	}

	realHash, err := getHashOfFile(filepath)
	if err != nil {
		return false, err
	}

	return hash == realHash, nil
}
