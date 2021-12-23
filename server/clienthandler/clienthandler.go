package clienthandler

import (
	"fmt"
	"io"

	"github.com/NikolayDPaev/CentralisedVersionControl/server/fileIO"
	"github.com/NikolayDPaev/CentralisedVersionControl/server/netIO"
)

const (
	GET_COMMIT_LIST_CODE = 0
	DOWNLOAD_COMMIT      = 1
	UPLOAD_COMMIT        = 2
)

func Communication(reader io.Reader, writer io.Writer) error {
	opCode, err := netIO.ReceiveUint32(reader)
	if err != nil {
		return fmt.Errorf("could not receive opcode: %w", err)
	}

	switch opCode {
	case GET_COMMIT_LIST_CODE:
		err = sendCommitList(writer)
	}

	return err
}

func sendCommitList(writer io.Writer) error {
	metadataList := fileIO.CommitList()
	err := netIO.SendUint32(uint32(len(metadataList)), writer)
	if err != nil {
		return fmt.Errorf("could not send metadata list length: %w", err)
	}

	for _, entry := range metadataList {
		err := netIO.SendString(entry.String(), writer)
		if err != nil {
			return fmt.Errorf("could not send metadata entry: %w", err)
		}
	}
	return nil
}
