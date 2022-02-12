package storage

import (
	"errors"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/NikolayDPaev/CentralisedVersionControl/netio"
	"github.com/NikolayDPaev/CentralisedVersionControl/server/servercommit"
)

// Predicate that checks if file with the specified filePath exists.
func fileExists(filePath string) (bool, error) {
	if _, err := os.Stat(filePath); err == nil {
		return true, nil

	} else if errors.Is(err, os.ErrNotExist) {
		return false, nil

	} else {
		return false, err
	}
}

// Returns a slice with the string representation of all commits on the disk,
// sorted in the order of creation.
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
	i := 0
	for _, v := range files {
		commit, err := s.OpenCommit(v.Name())
		if err != nil {
			log.Println(err.Error())
			continue
		}
		result[i] = commit.String()
		i++
	}
	return result
}

// Deserializes a new commit with this commit id from the disk.
func (s *FileStorage) OpenCommit(commitId string) (*servercommit.Commit, error) {
	file, err := os.Open("commits/" + commitId)
	if err != nil {
		return nil, fmt.Errorf("cannot open commit %s: %w", commitId, err)
	}
	fileComm := netio.NewCommunicator(0, file, file)
	return servercommit.NewCommitFrom(commitId, fileComm)
}

// Serializes the provided commit record to the disk.
func (s *FileStorage) SaveCommit(commit *servercommit.Commit) error {
	if err := os.MkdirAll("commits", 0777); err != nil {
		return fmt.Errorf("cannot create commit folder: %w", err)
	}

	file, err := os.Create("commits/" + commit.Id)
	if err != nil {
		return fmt.Errorf("cannot create commit file %s: %w", commit.Id, err)
	}
	defer file.Close()

	comm := netio.NewCommunicator(100, file, file)
	if err := commit.WriteTo(comm); err != nil {
		return fmt.Errorf("error saving commit %s: %w", commit.String(), err)
	}
	return nil
}

// Returns the size of the commit file on the disk.
func (s *FileStorage) CommitSize(commitId string) (int64, error) {
	fileInfo, err := os.Stat("commits/" + commitId)
	if err != nil {
		return 0, fmt.Errorf("cannot get commit %s file info: %w", commitId, err)
	}

	return fileInfo.Size(), nil
}

// Predicate that checks if there is a commit with the specified commit id on the disk.
func (s *FileStorage) CommitExists(commitId string) (bool, error) {
	b, err := fileExists("commits/" + commitId)
	if err != nil {
		return false, err
	}
	return b, nil
}
