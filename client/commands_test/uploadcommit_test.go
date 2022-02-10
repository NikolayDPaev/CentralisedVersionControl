package commands_test

import (
	"testing"

	"github.com/NikolayDPaev/CentralisedVersionControl/client/clientcommit"
	"github.com/NikolayDPaev/CentralisedVersionControl/client/commands"
	"github.com/NikolayDPaev/CentralisedVersionControl/client/fileio/fileiofakes"
	"github.com/NikolayDPaev/CentralisedVersionControl/netio/netiofakes"
)

type uploadTestCase struct {
	expectedCommit clientcommit.Commit
	missingBlobIds []string
}

func slicesFromSortedSliceEntries(fileSortedSlice []clientcommit.CommitEntry) ([]string, []string) {
	paths := make([]string, len(fileSortedSlice))
	blobIds := make([]string, len(fileSortedSlice))
	i := 0
	for _, entry := range fileSortedSlice {
		blobIds[i] = entry.Hash
		paths[i] = entry.Path
		i++
	}
	return paths, blobIds
}

const DW_OPCODE = 2

func TestUploadcommit(t *testing.T) {
	testcases := []uploadTestCase{
		{clientcommit.Commit{
			Message: "some message",
			Creator: "user",
			FileSortedSlice: []clientcommit.CommitEntry{
				{Hash: "123", Path: "a/path/to/file"},
				{Hash: "234", Path: "i/am/getting/bored"},
				{Hash: "456", Path: "another/path/to/other/file"},
				{Hash: "789", Path: "yet/again/a/path"}}},
			[]string{"123", "789"}},
		{clientcommit.Commit{
			Message:         "empty commit",
			Creator:         "user",
			FileSortedSlice: []clientcommit.CommitEntry{}},
			[]string{}},
	}
	for _, testcase := range testcases {
		commFake := &netiofakes.FakeCommunicator{}
		fileFake := &fileiofakes.FakeLocalcopy{}
		paths, blobIds := slicesFromSortedSliceEntries(testcase.expectedCommit.FileSortedSlice)

		fileFake.GetPathsOfAllFilesReturnsOnCall(0, paths, nil)
		for i, blob := range blobIds {
			fileFake.GetHashOfFileReturnsOnCall(i, blob, nil)
		}

		//sending missing blobs
		commFake.RecvStringSliceReturnsOnCall(0, testcase.missingBlobIds, nil)

		uploadCommit := commands.NewUpload(commFake, fileFake, DW_OPCODE)
		if err := uploadCommit.UploadCommit(testcase.expectedCommit.Message, testcase.expectedCommit.Creator); err != nil {
			t.Errorf("Error catched from Upload commit: %v", err)
		}

		if actual := commFake.SendVarIntArgsForCall(0); actual != DW_OPCODE {
			t.Errorf("Send var int called with wrong arg, when sending opcode: Expected: %d, actual: %d", DW_OPCODE, actual)
		}

		//creating commit
		for i, path := range paths {
			if actual := fileFake.GetHashOfFileArgsForCall(i); actual != path {
				t.Errorf("Get hash of file called with wrong args: Expected: %s, actual: %s", path, actual)
			}
		}

		//sending commit
		commitHash := testcase.expectedCommit.Md5Hash()
		if actual := commFake.SendStringArgsForCall(0); actual != commitHash {
			t.Errorf("Send string called with wrong args when sending commit hash: Expected: %s, actual: %s", commitHash, actual)
		}

		if actual := commFake.SendStringArgsForCall(1); actual != testcase.expectedCommit.Message {
			t.Errorf("Send string called with wrong args when sending commit message: Expected: %s, actual: %s", testcase.expectedCommit.Message, actual)
		}

		if actual := commFake.SendStringArgsForCall(2); actual != testcase.expectedCommit.Creator {
			t.Errorf("Send string called with wrong args when sending commit creator: Expected: %s, actual: %s", testcase.expectedCommit.Creator, actual)
		}

		commitTree := testcase.expectedCommit.GetTree()
		if actual := commFake.SendStringArgsForCall(3); actual != commitTree {
			t.Errorf("Send string called with wrong args when sending commit tree: Expected: %s, actual: %s", commitTree, actual)
		}

		//sending missing blobs
		for i, blobId := range testcase.missingBlobIds {
			if actual := commFake.SendStringArgsForCall(4 + i); actual != blobId {
				t.Errorf("Send string called with wrong args when sending blob id: Expected: %s, actual: %s", blobId, actual)
			}
			expectedPath, _ := testcase.expectedCommit.GetBlobPath(blobId)
			if actual, _ := fileFake.SendBlobArgsForCall(i); actual != expectedPath {
				t.Errorf("Send string called with wrong args when sending blob id: Expected: %s, actual: %s", expectedPath, actual)
			}
		}

	}
}
