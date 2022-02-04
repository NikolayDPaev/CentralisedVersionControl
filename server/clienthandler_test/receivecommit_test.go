package clienthandler_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/NikolayDPaev/CentralisedVersionControl/netio/netiofakes"
	"github.com/NikolayDPaev/CentralisedVersionControl/server/clienthandler"
	"github.com/NikolayDPaev/CentralisedVersionControl/server/storage/storagefakes"
)

type testCaseData struct {
	stringsToReceive []string
	missingBlobs     []string
}

func TestReceiveCommitHandle(t *testing.T) {
	var testCases = []testCaseData{
		{[]string{"12345", "message", "creator",
			"123 /some/path/to/blob\n456 /different/path\n789 /another/path\n101 /and/another", // commit data
			"789", "456"}, // sended blobs
			[]string{"456", "789"}}, // missing blobs

		{[]string{"12345", "empty commit", "creator",
			""}, // empty commit data
			[]string{}},
	}
	for _, testCase := range testCases {
		netFake := &netiofakes.FakeCommunicator{}
		fileFake := &storagefakes.FakeStorage{}

		for i, strings := range testCase.stringsToReceive {
			netFake.RecvStringReturnsOnCall(i, strings, nil)
		}
		fileFake.BlobExistsStub = func(s string) (bool, error) {
			for _, blobid := range testCase.missingBlobs {
				if blobid == s {
					return false, nil
				}
			}
			return true, nil
		}
		commitList := clienthandler.NewReceiveCommit(netFake, fileFake)
		if err := commitList.Handle(); err != nil {
			t.Errorf("Error catched from commitList.Handle(): %s", err)
		}

		if len(testCase.missingBlobs) > 0 {
			actualMissingBlobs := netFake.SendStringSliceArgsForCall(0)
			if !reflect.DeepEqual(testCase.missingBlobs, actualMissingBlobs) {
				t.Errorf("Wrong missing blobs sent. Expected: %s, actual: %s", testCase.missingBlobs, actualMissingBlobs)
			}
		}

		if len(testCase.stringsToReceive) > 4 {
			for i, blobId := range testCase.stringsToReceive[4:] {
				actualMissingBlobId, _ := fileFake.SaveBlobArgsForCall(i)
				if blobId != actualMissingBlobId {
					t.Errorf("SaveBlob called with wrong blob id. Expected: %s, actual: %s", blobId, actualMissingBlobId)
				}
			}
		}

		actualCommit := fileFake.SaveCommitArgsForCall(0)
		if strings.ReplaceAll(actualCommit.String(), "\"", "") != strings.Join(testCase.stringsToReceive[:3], " ") {
			t.Errorf("Save commit called with wrong commit")
		}
	}
}
