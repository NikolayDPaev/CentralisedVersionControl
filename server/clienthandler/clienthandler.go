package clienthandler

import (
	"fmt"

	"github.com/NikolayDPaev/CentralisedVersionControl/server/fileIO"
	"github.com/NikolayDPaev/CentralisedVersionControl/server/netIO"
)

const (
	GET_COMMIT_LIST = 0
	DOWNLOAD_COMMIT = 1
	UPLOAD_COMMIT   = 2
	EMPTY_REQUEST   = 3
)

type Clienthandler interface {
	Handle() error
}

func NewHandler(comm netIO.Communicator, storage fileIO.Storage) (Clienthandler, error) {
	opCode, err := comm.ReceiveVarInt()
	if err != nil {
		return nil, fmt.Errorf("could not receive opcode:\n%w", err)
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
