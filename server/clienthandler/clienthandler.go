package clienthandler

import (
	"fmt"

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

func NewHandler(comm netIO.Communicator) (Clienthandler, error) {
	opCode, err := comm.ReceiveVarInt()
	if err != nil {
		return nil, fmt.Errorf("could not receive opcode:\n%w", err)
	}

	switch opCode {
	case GET_COMMIT_LIST:
		return &CommitList{comm}, nil
	case UPLOAD_COMMIT:
		return &ReceiveCommit{comm}, nil
	case DOWNLOAD_COMMIT:
		return &SendCommit{comm}, nil
	case EMPTY_REQUEST:
		return nil, nil
	}

	return nil, fmt.Errorf("invalid opcode %d", opCode)
}
