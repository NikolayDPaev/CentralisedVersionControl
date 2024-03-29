// Defines the operations that the server provides
package clienthandler

import (
	"fmt"

	"github.com/NikolayDPaev/CentralisedVersionControl/netio"
	"github.com/NikolayDPaev/CentralisedVersionControl/server/storage"
)

// Operation codes
const (
	GET_COMMIT_LIST = 0
	DOWNLOAD_COMMIT = 1
	UPLOAD_COMMIT   = 2
	EMPTY_REQUEST   = 3
)

type Clienthandler interface {
	Handle() error
}

// NewHandler returns clienthandler implementation based on the opcode received from the client
func NewHandler(comm netio.Communicator, storage storage.Storage) (Clienthandler, error) {
	opCode, err := comm.RecvVarInt()
	if err != nil {
		return nil, fmt.Errorf("could not receive opcode: %w", err)
	}

	switch opCode {
	case GET_COMMIT_LIST:
		return NewCommitList(comm, storage), nil
	case UPLOAD_COMMIT:
		return NewReceiveCommit(comm, storage), nil
	case DOWNLOAD_COMMIT:
		return NewSendCommit(comm, storage), nil
	case EMPTY_REQUEST:
		return nil, nil
	}

	return nil, fmt.Errorf("invalid opcode %d", opCode)
}
