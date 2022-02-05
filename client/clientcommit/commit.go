package clientcommit

import (
	"crypto/md5"
	"encoding/base32"
	"errors"
	"fmt"
	"strings"

	"github.com/NikolayDPaev/CentralisedVersionControl/client/fileio"
)

type Commit struct {
	Message string
	Creator string
	FileMap map[string]string
}

func GetMap(tree string) (map[string]string, error) {
	if len(tree) == 0 {
		return nil, nil
	}
	lines := strings.Split(tree, "\n")

	fileMap := make(map[string]string, len(lines))

	for _, line := range lines {
		elements := strings.Split(line, " ")
		if len(elements) != 2 {
			return nil, errors.New("corrupt commit tree string")
		}
		fileMap[elements[0]] = elements[1]
	}

	return fileMap, nil
}

func (c *Commit) GetBlobPath(blobId string) (string, error) {
	path, ok := c.FileMap[blobId]
	if !ok {
		return "", errors.New("missing blobId in filemap")
	}
	return path, nil
}

func (c *Commit) GetTree() string {
	var sb strings.Builder
	for id, path := range c.FileMap {
		sb.WriteString(id + " " + path + "\n")
	}
	str := sb.String()
	if sb.Len() > 1 {
		return str[:len(str)-1] // removing trailing endline
	}
	return str
}

func (c *Commit) GetMissingFiles(localcpy fileio.Localcopy) (map[string]string, error) {
	missingFileMap := make(map[string]string, len(c.FileMap)/2)

	for blobId, path := range c.FileMap {
		exists, err := localcpy.FileWithHashExists(path, blobId)
		if err != nil {
			return nil, err
		}

		if !exists {
			missingFileMap[blobId] = path
		}
	}

	return missingFileMap, nil
}

func (c *Commit) Md5Hash() string {
	hash := md5.Sum([]byte(fmt.Sprintf("%v", c)))
	str := base32.StdEncoding.EncodeToString(hash[:])
	return str
}

func (c *Commit) GetSetOfPaths() map[string]struct{} {
	set := make(map[string]struct{}, len(c.FileMap))

	for _, values := range c.FileMap {
		set[values] = struct{}{}
	}

	return set
}
