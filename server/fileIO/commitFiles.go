package fileIO

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"sort"

	"github.com/NikolayDPaev/CentralisedVersionControl/server/commit"
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

func extractMetadata(fileInfo fs.FileInfo) (*commit.Metadata, error) {
	file, err := os.Open("commits/" + fileInfo.Name())
	if err != nil {
		return nil, err
	}
	defer file.Close()

	result, err := commit.ReadMetadata(file, fileInfo.Name())
	if err != nil {
		return nil, err
	}
	return result, nil
}

func CommitList() []*commit.Metadata {
	f, err := os.Open("./commits")
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	defer f.Close()

	files, err := f.Readdir(0)
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].ModTime().Before(files[j].ModTime())
	})

	result := make([]*commit.Metadata, len(files))
	for i, v := range files {
		result[i], err = extractMetadata(v)
		if err != nil {
			continue
		}
	}
	return result
}

func OpenCommit(commitId string) (*os.File, error) {
	file, err := os.Open("commits/" + commitId)
	if err != nil {
		return nil, fmt.Errorf("cannot open commit %s: %w", commitId, err)
	}
	return file, nil
}

func NewCommit(commitId string) (*os.File, error) {
	file, err := os.Create("commits/" + commitId)
	if err != nil {
		return nil, fmt.Errorf("cannot create commit file %s: %w", commitId, err)
	}
	return file, nil
}

func CommitSize(commitId string) (int64, error) {
	fileInfo, err := os.Stat("commits/" + commitId)
	if err != nil {
		return 0, fmt.Errorf("cannot get commit %s file info: %w", commitId, err)
	}

	return fileInfo.Size(), nil
}

func CommitExists(commitId string) (bool, error) {
	b, err := fileExists("commits/" + commitId)
	if err != nil {
		return false, err
	}
	return b, nil
}
