package fileio

import "os"

type Localcopy interface {
	GetHashOfFile(filepath string) (string, error)
	GetPathsOfAllFiles() ([]string, error)
	FileWithHashExists(filepath string, hash string) (bool, error)
	FileSize(path string) (int64, error)
	CleanOtherFiles(commitFilesSet map[string]struct{}) error
	CompressToTempFile(source string) (*os.File, error)
	DecompressFile(dest string, sFile *os.File) error
}

type Localfiles struct{}
