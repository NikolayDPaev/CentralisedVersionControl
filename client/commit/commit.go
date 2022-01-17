package commit

import (
	"crypto/md5"
	"encoding/base32"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/NikolayDPaev/CentralisedVersionControl/client/fileIO"
	"github.com/NikolayDPaev/CentralisedVersionControl/client/netIO"
)

type Commit struct {
	message string
	creator string
	fileMap map[string]string
}

func (c *Commit) GetBlobPath(blobId string) (string, error) {
	path, ok := c.fileMap[blobId]
	if !ok {
		return "", errors.New("missing blobId in filemap")
	}
	return path, nil
}

func getTree(fileMap map[string]string) string {
	var sb strings.Builder
	for id, path := range fileMap {
		sb.WriteString(id + " " + path + "\n")
	}
	str := sb.String()
	if sb.Len() > 1 {
		return str[:len(str)-1] // removing trailing endline
	}
	return str
}

func (c *Commit) Send(writer io.Writer) error {
	err := netIO.SendString(c.message, writer)
	if err != nil {
		return fmt.Errorf("cannot send commit message:\n%w", err)
	}

	err = netIO.SendString(c.creator, writer)
	if err != nil {
		return fmt.Errorf("cannot send commit creator:\n%w", err)
	}

	err = netIO.SendString(getTree(c.fileMap), writer)
	if err != nil {
		return fmt.Errorf("error sending commit tree:\n%w", err)
	}

	return nil
}
func (c *Commit) GetMissingFiles() (map[string]string, error) {
	missingFileMap := make(map[string]string, len(c.fileMap)/2)

	for blobId, path := range c.fileMap {
		exists, err := fileIO.FileWithHashExists(path, blobId)
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
	set := make(map[string]struct{}, len(c.fileMap))

	for _, values := range c.fileMap {
		set[values] = struct{}{}
	}

	return set
}
