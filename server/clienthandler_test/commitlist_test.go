package clienthandler_test

import (
	"reflect"
	"testing"

	"github.com/NikolayDPaev/CentralisedVersionControl/netio/netiofakes"
	"github.com/NikolayDPaev/CentralisedVersionControl/server/clienthandler"
	"github.com/NikolayDPaev/CentralisedVersionControl/server/storage/storagefakes"
)

func TestCommitListHandle(t *testing.T) {
	testCases := [][]string{
		{},
		{"some", "strings", "representing", "commits"},
	}
	for _, testCase := range testCases {
		netFake := &netiofakes.FakeCommunicator{}
		fileFake := &storagefakes.FakeStorage{}

		fileFake.CommitListReturns(testCase)

		commitList := clienthandler.NewCommitList(netFake, fileFake)
		commitList.Handle()

		actualArgs := netFake.SendStringSliceArgsForCall(0)
		if !reflect.DeepEqual(testCase, actualArgs) {
			t.Errorf("Send string called with wrong args. Expected: %s, actual: %s", testCase, actualArgs)
		}
	}
}
