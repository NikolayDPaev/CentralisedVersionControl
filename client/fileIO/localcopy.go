package fileio

import (
	"github.com/NikolayDPaev/CentralisedVersionControl/netio"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Localcopy

type Localcopy interface {
	GetHashOfFile(filepath string) (string, error)
	GetPathsOfAllFiles() ([]string, error)
	FileWithHashExists(filepath string, hash string) (bool, error)
	FileSize(path string) (int64, error)
	CleanOtherFiles(commitFilesSet map[string]struct{}) error
	ReceiveBlob(filepath string, comm netio.Communicator) error
	SendBlob(filepath string, comm netio.Communicator) error
}

type Localfiles struct {
	fileExceptions map[string]struct{}
}

func NewLocalfiles(fileExceptions map[string]struct{}) *Localfiles {
	return &Localfiles{fileExceptions: fileExceptions}
}
