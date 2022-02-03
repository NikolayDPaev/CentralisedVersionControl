package clientcommit

import (
	"errors"
	"fmt"
	"strings"

	"github.com/NikolayDPaev/CentralisedVersionControl/client/fileio"
	"github.com/NikolayDPaev/CentralisedVersionControl/netio"
)

func readMetadata(comm netio.Communicator) (string, string, error) {
	message, err := comm.RecvString()
	if err != nil {
		return "", "", err
	}

	creator, err := comm.RecvString()
	if err != nil {
		return "", "", err
	}

	return message, creator, nil
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

func ReadCommit(id string, comm netio.Communicator) (*Commit, error) {
	message, creator, err := readMetadata(comm)
	if err != nil {
		return nil, fmt.Errorf("cannot read metadata of commit:\n%w", err)
	}

	strTree, err := comm.RecvString()
	if err != nil {
		return nil, fmt.Errorf("cannot read tree string of commit:\n%w", err)
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

func CreateCommit(message, creator string, localcpy fileio.Localcopy) (*Commit, error) {
	paths, err := localcpy.GetPathsOfAllFiles()
	if err != nil {
		return nil, fmt.Errorf("error getting filenames for creating commit:\n%w", err)
	}

	fileMap := make(map[string]string, len(paths))
	for _, path := range paths {
		hash, err := localcpy.GetHashOfFile(path)
		if err != nil {
			return nil, err
		}
		fileMap[hash] = path
	}

	return &Commit{message: message, creator: creator, fileMap: fileMap}, nil
}
