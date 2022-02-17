package fileio

import (
	"github.com/NikolayDPaev/CentralisedVersionControl/netio"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Localcopy

// Localcopy defines all operations with the local files
type Localcopy interface {
	// GetHashOfFile returns Md5Sum of the file with the provided filepath
	GetHashOfFile(filepath string) (string, error)

	// FileWithHashExists checks if the file with the provided path has the provided hash.
	// If file does not exist returns error.
	FileWithHashExists(filepath string, hash string) (bool, error)

	// GetPathsOfAllFiles returns slice with paths to all files in all directories.
	GetPathsOfAllFiles() ([]string, error)

	// CleanOtherFiles deletes all files that are not part of the provided commitFilesSet.
	CleanOtherFiles(commitFilesSet map[string]struct{}) error

	// ReceiveBlob receives the file with the provided path via the communicator.
	ReceiveBlob(filepath string, comm netio.Communicator) error

	// SendBlob sends the file with the provided path via the communicator.
	SendBlob(filepath string, comm netio.Communicator) error
}

// Localfiles is an implementation of the Localcopy interface
// Contains fileExceptions field that marks all files that should not be
// included in commits or deleted
type Localfiles struct {
	ignoredFiles map[string]struct{}
}

func NewLocalfiles(ignoredFiles map[string]struct{}) *Localfiles {
	return &Localfiles{ignoredFiles: ignoredFiles}
}
