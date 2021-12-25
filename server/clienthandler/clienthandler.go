package clienthandler

import (
	"fmt"
	"io"

	"github.com/NikolayDPaev/CentralisedVersionControl/server/netIO"
)

const (
	GET_COMMIT_LIST = 0
	DOWNLOAD_COMMIT = 1
	UPLOAD_COMMIT   = 2
	EMPTY_REQUEST   = 3
)

func Communication(reader io.Reader, writer io.Writer) error {
	opCode, err := netIO.ReceiveVarInt(reader)
	if err != nil {
		return fmt.Errorf("could not receive opcode: %w", err)
	}

	switch opCode {
	case GET_COMMIT_LIST:
		err = sendCommitList(writer)
	case UPLOAD_COMMIT:
		err = receiveCommit(reader, writer)
	case DOWNLOAD_COMMIT:
		err = sendCommit(reader, writer)
	case EMPTY_REQUEST:
		return nil
	}

	return err
}
