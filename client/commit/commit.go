package commit

import (
	"bufio"
	"crypto/md5"
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

func getMap(tree string) (map[string]string, error) {
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

func readMetadata(reader io.Reader) (string, string, error) {
	scanner := bufio.NewScanner(reader)
	scanner.Scan()
	message := scanner.Text()
	scanner.Scan()
	creator := scanner.Text()

	if err := scanner.Err(); err != nil {
		return "", "", err
	}

	return message, creator, nil
}

func ReadCommit(reader io.Reader) (*Commit, error) {
	id, err := netIO.ReceiveString(reader)
	if err != nil {
		return nil, fmt.Errorf("cannot read id of commit: %w", err)
	}

	message, creator, err := readMetadata(reader)
	if err != nil {
		return nil, fmt.Errorf("cannot read metadata of commit: %w", err)
	}

	strTree, err := netIO.ReceiveString(reader)
	if err != nil {
		return nil, fmt.Errorf("cannot read tree string of commit: %w", err)
	}

	fileMap, err := getMap(strTree)
	if err != nil {
		return nil, err
	}

	commit := &Commit{message, creator, fileMap}
	if id != commit.Md5Hash() {
		return nil, errors.New("mismatched hash values")
	}
	return commit, nil
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

	return string(hash[:])
}
