package fileIO

import (
	"crypto/md5"
	"encoding/base32"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
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

func GetHashOfFile(filepath string) (string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "", fmt.Errorf("error opening %s:\n%w", filepath, err)
	}
	defer file.Close()

	hash := md5.New()
	_, err = io.Copy(hash, file)

	if err != nil {
		return "", fmt.Errorf("error getting hash of %s:\n%w", filepath, err)
	}
	str := base32.StdEncoding.EncodeToString(hash.Sum(nil))
	return str, nil
}

func FileWithHashExists(filepath string, hash string) (bool, error) {
	fileExists, err := fileExists(filepath)
	if err != nil {
		return false, fmt.Errorf("error checking if file exists:\n%w", err)
	}
	if !fileExists {
		return false, nil
	}

	realHash, err := GetHashOfFile(filepath)
	if err != nil {
		return false, err
	}

	return hash == realHash, nil
}

func FileSize(path string) (int64, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return 0, fmt.Errorf("cannot get file %s file info:\n%w", path, err)
	}

	return fileInfo.Size(), nil
}

func GetPathsOfAllFiles() ([]string, error) {
	var paths []string
	var stack []string

	stack = append(stack, ".") // windows?
	for len(stack) > 0 {
		n := len(stack) - 1
		curDir := stack[n] // top
		stack = stack[:n]  // pop

		files, err := ioutil.ReadDir(curDir)
		if err != nil {
			return nil, fmt.Errorf("error scanning directory %s:\n%w", curDir, err)
		}

		for _, file := range files {
			if file.IsDir() {
				stack = append(stack, curDir+"/"+file.Name())
			} else {
				paths = append(paths, curDir+"/"+file.Name())
			}
		}
	}

	return paths, nil
}
