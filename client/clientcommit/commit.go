package clientcommit

import (
	"crypto/md5"
	"encoding/base32"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/NikolayDPaev/CentralisedVersionControl/client/fileio"
)

type CommitEntry struct {
	Hash string
	Path string
}

type Commit struct {
	Message         string
	Creator         string
	FileSortedSlice []CommitEntry
}

func GetSortedSlice(tree string) ([]CommitEntry, error) {
	if len(tree) == 0 {
		return nil, nil
	}
	lines := strings.Split(tree, "\n")

	fileSortedSlice := make([]CommitEntry, len(lines))
	for i, line := range lines {
		elements := strings.Split(line, " ")
		if len(elements) != 2 {
			return nil, errors.New("corrupt commit tree string")
		}
		fileSortedSlice[i] = CommitEntry{Hash: elements[0], Path: elements[1]}
	}

	sort.Slice(fileSortedSlice, func(i, j int) bool {
		return fileSortedSlice[i].Hash < fileSortedSlice[j].Hash
	})

	return fileSortedSlice, nil
}

func (c *Commit) GetBlobPath(blobId string) (string, error) {
	start := 0
	end := len(c.FileSortedSlice) - 1
	for start <= end {
		mid := (start + end) / 2

		if c.FileSortedSlice[mid].Hash == blobId {
			return c.FileSortedSlice[mid].Path, nil
		} else if c.FileSortedSlice[mid].Hash < blobId {
			start = mid + 1
		} else if c.FileSortedSlice[mid].Hash > blobId {
			end = mid - 1
		}
	}
	return "", errors.New("missing blobId in file slice")
}

func (c *Commit) GetTree() string {
	var sb strings.Builder
	for _, entry := range c.FileSortedSlice {
		sb.WriteString(entry.Hash + " " + entry.Path + "\n")
	}
	str := sb.String()
	if sb.Len() > 1 {
		return str[:len(str)-1] // removing trailing endline
	}
	return str
}

func (c *Commit) GetMissingFiles(localcpy fileio.Localcopy) (map[string]string, error) {
	missingFileMap := make(map[string]string, len(c.FileSortedSlice)/2)

	for _, entry := range c.FileSortedSlice {
		exists, err := localcpy.FileWithHashExists(entry.Path, entry.Hash)
		if err != nil {
			return nil, err
		}

		if !exists {
			missingFileMap[entry.Hash] = entry.Path
		}
	}

	return missingFileMap, nil
}

func (c *Commit) Md5Hash() string {
	hash := md5.Sum([]byte(fmt.Sprintf("%v", c)))
	str := base32.StdEncoding.EncodeToString(hash[:])

	return strings.ReplaceAll(str, "=", "")
}

func (c *Commit) GetSetOfPaths() map[string]struct{} {
	set := make(map[string]struct{}, len(c.FileSortedSlice))

	for _, entry := range c.FileSortedSlice {
		set[entry.Path] = struct{}{}
	}

	return set
}
