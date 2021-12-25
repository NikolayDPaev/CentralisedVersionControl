package commands

import (
	"fmt"
	"io"

	"github.com/NikolayDPaev/CentralisedVersionControl/client/netIO"
)

const (
	GET_COMMIT_LIST = 0
	DOWNLOAD_COMMIT = 1
	UPLOAD_COMMIT   = 2
	EMPTY_REQUEST   = 3
)

func GetCommitList(reader io.Reader) ([]string, error) {
	err := netIO.SendVarInt(GET_COMMIT_LIST, reader)
	if err != nil {
		return nil, fmt.Errorf("cannot send op code: %w", err)
	}

	commitList, err := netIO.ReceiveStringSlice(reader)
	if err != nil {
		return nil, fmt.Errorf("error receiving commit list: %w", err)
	}
	return commitList, nil
}
