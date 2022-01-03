package clienthandler

import (
	"fmt"
	"io"

	"github.com/NikolayDPaev/CentralisedVersionControl/server/fileIO"
	"github.com/NikolayDPaev/CentralisedVersionControl/server/netIO"
)

func sendCommitList(writer io.Writer) error {
	metadataList := fileIO.CommitList()
	err := netIO.SendVarInt(int64(len(metadataList)), writer)
	if err != nil {
		return fmt.Errorf("could not send metadata list length:\n%w", err)
	}

	for _, entry := range metadataList {
		err := netIO.SendString(entry, writer)
		if err != nil {
			return fmt.Errorf("could not send metadata entry:\n%w", err)
		}
	}
	return nil
}
