package clientcommit

import (
	"fmt"

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

func ReadCommit(id string, comm netio.Communicator) (*Commit, error) {
	receivedId, err := comm.RecvString()
	if err != nil || id != receivedId {
		return nil, fmt.Errorf("error receiving id of commit:\n%w", err)
	}

	message, creator, err := readMetadata(comm)
	if err != nil {
		return nil, fmt.Errorf("cannot read metadata of commit:\n%w", err)
	}

	strTree, err := comm.RecvString()
	if err != nil {
		return nil, fmt.Errorf("cannot read tree string of commit:\n%w", err)
	}

	fileMap, err := GetMap(strTree)
	if err != nil {
		return nil, err
	}

	commit := &Commit{message, creator, fileMap}
	commitHash := commit.Md5Hash()
	if id != commitHash {
		return nil, fmt.Errorf("mismatched hash values: expected: %s, actual: %s", id, commitHash)
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

	return &Commit{Message: message, Creator: creator, FileMap: fileMap}, nil
}
