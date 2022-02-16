package fileio

import (
	"github.com/NikolayDPaev/CentralisedVersionControl/netio"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Localcopy

// Interface that defines all operations with the local files
type Localcopy interface {
	// Returns Md5Sum of the file with the provided filepath
	GetHashOfFile(filepath string) (string, error)

	// Checks if the file with the provided path has the provided hash.
	// If file does not exist returns error.
	FileWithHashExists(filepath string, hash string) (bool, error)

	// Returns slice with paths to all files in all directories.
	GetPathsOfAllFiles() ([]string, error)

	// Deletes all files that are not part of the provided commitFilesSet.
	CleanOtherFiles(commitFilesSet map[string]struct{}) error

	// Receives the file with the provided path via the communicator.
	ReceiveBlob(filepath string, comm netio.Communicator) error

	// Sends the file with the provided path via the communicator.
	SendBlob(filepath string, comm netio.Communicator) error
}

// Implementation of the Localcopy interface
// Contains fileExceptions field that marks all files that should not be
// included in commits or deleted
type Localfiles struct {
	ignoredFiles map[string]struct{}
}

func NewLocalfiles(ignoredFiles map[string]struct{}) *Localfiles {
	return &Localfiles{ignoredFiles: ignoredFiles}
}
