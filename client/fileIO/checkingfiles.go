package fileio

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
		return "", fmt.Errorf("error opening %s: %w", filepath, err)
	}
	defer file.Close()

	hash := md5.New()
	_, err = io.Copy(hash, file)

	if err != nil {
		return "", fmt.Errorf("error getting hash of %s: %w", filepath, err)
	}
	str := base32.StdEncoding.EncodeToString(hash.Sum(nil))
	return str, nil
}

func FileWithHashExists(filepath string, hash string) (bool, error) {
	fileExists, err := fileExists(filepath)
	if err != nil {
		return false, fmt.Errorf("error checking if file exists: %w", err)
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
		return 0, fmt.Errorf("cannot get file %s file info: %w", path, err)
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
			return nil, fmt.Errorf("error scanning directory %s: %w", curDir, err)
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

func removeDirIfEmpty(dir string) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("error scanning directory %s: %w", dir, err)
	}
	if len(files) == 0 {
		if err := os.Remove(dir); err != nil {
			return fmt.Errorf("error deleting empty directory %s: %w", dir, err)
		}
	}
	return nil
}

func deleteFileIfNotInSet(file string, filesSet map[string]struct{}) error {
	if _, ok := filesSet[file]; !ok {
		if err := os.Remove(file); err != nil {
			return fmt.Errorf("error deleting file %s: %w", file, err)
		}
	}
	return nil
}

func CleanOtherFiles(commitFilesSet map[string]struct{}) error {
	var stack []string
	stack = append(stack, ".")
	for len(stack) > 0 {
		n := len(stack) - 1
		curDir := stack[n] // top
		stack = stack[:n]  // pop

		files, err := ioutil.ReadDir(curDir)
		if err != nil {
			return fmt.Errorf("error scanning directory %s: %w", curDir, err)
		}

		for _, file := range files {
			if file.IsDir() {
				stack = append(stack, curDir+"/"+file.Name())
			} else {
				file := curDir + "/" + file.Name()
				if err := deleteFileIfNotInSet(file, commitFilesSet); err != nil {
					return err
				}
			}
		}
		if err := removeDirIfEmpty(curDir); err != nil {
			return err
		}
	}
	return nil
}
