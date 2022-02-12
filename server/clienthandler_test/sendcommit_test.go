package clienthandler_test

import (
	"testing"

	"github.com/NikolayDPaev/CentralisedVersionControl/netio/netiofakes"
	"github.com/NikolayDPaev/CentralisedVersionControl/server/clienthandler"
	"github.com/NikolayDPaev/CentralisedVersionControl/server/servercommit"
	"github.com/NikolayDPaev/CentralisedVersionControl/server/storage/storagefakes"
)

type testCase struct {
	commitId     string
	commit       servercommit.Commit
	blobsForSend []string
}

func TestSendCommitHandleUnavailableCommitId(t *testing.T) {
	commitId := "12345"

	commFake := &netiofakes.FakeCommunicator{}
	storeFake := &storagefakes.FakeStorage{}
	commFake.RecvStringReturnsOnCall(0, commitId, nil)
	storeFake.CommitExistsReturnsOnCall(0, false, nil)

	sendCommit := clienthandler.NewSendCommit(commFake, storeFake)
	if err := sendCommit.Handle(); err != nil {
		t.Errorf("Error when called handle: %v", err)
	}

	if okcode := commFake.SendVarIntArgsForCall(0); okcode != clienthandler.ERROR {
		t.Errorf("Send var int called with wrong arg when sending ok code. Expected %d, actual: %d", clienthandler.ERROR, okcode)
	}
}

func TestSendCommitHandle(t *testing.T) {
	testCases := []testCase{
		{"12345",
			servercommit.Commit{Id: "12345",
				Message: "random message",
				Creator: "creator",
				Tree:    "123 /some/path/to/blob\n456 /different/path\n789 /another/path\n101 /and/another"},
			[]string{"123", "789"}},
		{"12345",
			servercommit.Commit{Id: "12345",
				Message: "another message",
				Creator: "different creator",
				Tree:    "123 /some/path/to/blob\n456 /different/path\n789 /another/path\n101 /and/another"},
			[]string{}},
	}

	for _, testcase := range testCases {
		commFake := &netiofakes.FakeCommunicator{}
		storeFake := &storagefakes.FakeStorage{}

		commFake.RecvStringReturnsOnCall(0, testcase.commitId, nil)
		storeFake.CommitExistsReturnsOnCall(0, true, nil)
		storeFake.OpenCommitReturnsOnCall(0, &testcase.commit, nil)
		commFake.RecvStringSliceReturnsOnCall(0, testcase.blobsForSend, nil)

		sendCommit := clienthandler.NewSendCommit(commFake, storeFake)
		if err := sendCommit.Handle(); err != nil {
			t.Errorf("Error when called handle: %v", err)
		}

		if actual := storeFake.CommitExistsArgsForCall(0); testcase.commitId != actual {
			t.Errorf("Commit exists called with wrong arg. Expected %s, actual: %s", testcase.commitId, actual)
		}

		if okcode := commFake.SendVarIntArgsForCall(0); okcode != clienthandler.OK {
			t.Errorf("Send var int called with wrong arg when sending ok code. Expected %d, actual: %d", clienthandler.OK, okcode)
		}
		if testcase.commit.Id != commFake.SendStringArgsForCall(0) ||
			testcase.commit.Message != commFake.SendStringArgsForCall(1) ||
			testcase.commit.Creator != commFake.SendStringArgsForCall(2) ||
			testcase.commit.Tree != commFake.SendStringArgsForCall(3) {
			t.Errorf("Send string was not invoked correctly for sending of commit")
		}

		for i, blob := range testcase.blobsForSend {
			if actual, _ := storeFake.SendBlobArgsForCall(i); blob != actual {
				t.Errorf("Open blob called with wrong arg. Expected: %s, actual: %s", blob, actual)
			}
		}

		for i, blob := range testcase.blobsForSend {
			if actual := commFake.SendStringArgsForCall(i + 4); blob != actual {
				t.Errorf("Send string called with wrong arg when sending blobIds. Expected: %s, actual: %s", blob, actual)
			}
		}
	}
}
