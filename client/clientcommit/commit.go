// Package contains a data structure that represents a commit,
// along with some methods
package clientcommit

import (
	"crypto/md5"
	"encoding/base32"
	"errors"
	"fmt"
	"sort"
	"strings"
)

// Record, used to represent a file object - hash for id and path.
type CommitEntry struct {
	Hash string
	Path string
}

// Struct representing the client view of a commit.
// It has fields - Message and Creator and FileSortedSlice that
// keeps track of all files in the commit. The slice is sorted
// acording to the Hash in CommitEntry and searching for element
// based on the hash takes log n time.
//
// map[hash]path is not used because the order of elements
// in the map is undefined and that results in different
// Md5Sum of the whole commit.
// Since the Md5Sum is used as an ID of the commit the situation above
// may lead to having two identical commits with different ids.
type Commit struct {
	Message         string
	Creator         string
	FileSortedSlice []CommitEntry
}

// Function for transforming between the different representations of the files in commit
// Expects string that contains lines: "hash path"
// Returns sorted slice of CommitEntries
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

// Finds the path of a file by its hash.
// Uses binary search in the sorted slice representing the files.
// If there is no such blob with this hash - returns an error.
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

// Function for transforming between the different representations of the files in commit
// Returns string that contains lines: "hash path"
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

// Returns the Md5Sum of all fields in the struct
// It is used as an id of the commit
func (c *Commit) Md5Hash() string {
	hash := md5.Sum([]byte(fmt.Sprintf("%v", c)))
	str := base32.StdEncoding.EncodeToString(hash[:])

	return strings.ReplaceAll(str, "=", "")
}

// Returns a "set" - map[string]struct{} of all paths in the file slice
func (c *Commit) GetSetOfPaths() map[string]struct{} {
	set := make(map[string]struct{}, len(c.FileSortedSlice))

	for _, entry := range c.FileSortedSlice {
		set[entry.Path] = struct{}{}
	}

	return set
}
