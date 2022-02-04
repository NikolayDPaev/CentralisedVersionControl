package commands_test

import (
	"reflect"
	"testing"

	"github.com/NikolayDPaev/CentralisedVersionControl/client/client"
	"github.com/NikolayDPaev/CentralisedVersionControl/client/commands"
	"github.com/NikolayDPaev/CentralisedVersionControl/netio/netiofakes"
)

func TestGetCommitList(t *testing.T) {
	fakeComm := &netiofakes.FakeCommunicator{}

	expectedSlice := []string{"expected", "commitlist"}
	fakeComm.RecvStringSliceReturns(expectedSlice, nil)

	commitlist := commands.NewCommitList(fakeComm, client.GET_COMMIT_LIST)
	if actualSlice, _ := commitlist.GetCommitList(); !reflect.DeepEqual(expectedSlice, actualSlice) {
		t.Errorf("Commitlist returned wrong string. Expected: %s, actual: %s", expectedSlice, actualSlice)
	}

	if count := fakeComm.SendVarIntCallCount(); count != 2 {
		t.Errorf("SendVarInt different call count. Expected: %d, actual: %d", 2, count)
	}

	if opcode := fakeComm.SendVarIntArgsForCall(0); opcode != client.GET_COMMIT_LIST {
		t.Errorf("SendVarInt invoked with wrong args. Opcode is different. Expected: %d, actual: %d", client.GET_COMMIT_LIST, opcode)
	}
}
