package netio

import "io"

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Communicator

// Interface that defines methods for sending and receiving various data
type Communicator interface {
	SendVarInt(num int64) error
	RecvVarInt() (int64, error)
	SendString(str string) error
	RecvString() (string, error)
	SendStringSlice(slice []string) error
	RecvStringSlice() ([]string, error)
	SendFileData(fileReader io.Reader, fileLength int64) error
	RecvFileData(fileWriter io.Writer) error
}
