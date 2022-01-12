package netIO

import "io"

type Communicator interface {
	SendVarInt(num int64) error
	ReceiveVarInt() (int64, error)
	SendString(str string) error
	ReceiveString() (string, error)
	SendStringSlice(slice []string) error
	ReceiveStringSlice() ([]string, error)
	SendFileData(fileReader io.Reader, fileLength int64) error
	ReceiveFileData(fileWriter io.Writer) error
}
