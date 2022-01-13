package fileIO

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"sort"

	"github.com/NikolayDPaev/CentralisedVersionControl/server/commit"
	"github.com/NikolayDPaev/CentralisedVersionControl/server/netIO"
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

func extractCommitData(fileInfo fs.FileInfo) (string, error) { // !!!
	file, err := os.Open("commits/" + fileInfo.Name())
	if err != nil {
		return "", err
	}
	defer file.Close()

	comm := netIO.NewCommunicator(100, file, file)
	message, creator, err := commit.ReadCommitData(comm)
	if err != nil {
		return "", err
	}
	return fileInfo.Name() + " \"" + message + "\" " + creator, nil
}

func (s *FileStorage) CommitList() []string {
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

	result := make([]string, len(files))
	for i, v := range files {
		result[i], err = extractCommitData(v)
		if err != nil {
			continue
		}
	}
	return result
}

func (s *FileStorage) OpenCommit(commitId string) (StorageEntry, error) {
	file, err := os.Open("commits/" + commitId)
	if err != nil {
		return nil, fmt.Errorf("cannot open commit %s:\n%w", commitId, err)
	}
	return file, nil
}

func (s *FileStorage) NewCommit(commitId string) (StorageEntry, error) {
	if err := os.MkdirAll("commits", 0777); err != nil {
		return nil, fmt.Errorf("cannot create commit folder:\n%w", err)
	}

	file, err := os.Create("commits/" + commitId)
	if err != nil {
		return nil, fmt.Errorf("cannot create commit file %s:\n%w", commitId, err)
	}
	return file, nil
}

func (s *FileStorage) CommitSize(commitId string) (int64, error) {
	fileInfo, err := os.Stat("commits/" + commitId)
	if err != nil {
		return 0, fmt.Errorf("cannot get commit %s file info:\n%w", commitId, err)
	}

	return fileInfo.Size(), nil
}

func (s *FileStorage) CommitExists(commitId string) (bool, error) {
	b, err := fileExists("commits/" + commitId)
	if err != nil {
		return false, err
	}
	return b, nil
}
