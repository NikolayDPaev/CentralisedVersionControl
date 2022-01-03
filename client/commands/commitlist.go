package commands

import (
	"fmt"
	"io"

	"github.com/NikolayDPaev/CentralisedVersionControl/client/netIO"
)

const GET_COMMIT_LIST = 0

func GetCommitList(reader io.Reader, writer io.Writer) ([]string, error) {
	err := netIO.SendVarInt(GET_COMMIT_LIST, writer)
	if err != nil {
		return nil, fmt.Errorf("cannot send op code:\n%w", err)
	}

	commitList, err := netIO.ReceiveStringSlice(reader)
	if err != nil {
		return nil, fmt.Errorf("error receiving commit list:\n%w", err)
	}
	return commitList, nil
}
