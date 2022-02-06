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

func filesFromMap(fileMap map[string]string) ([]string, []string) {
	paths := make([]string, len(fileMap))
	blobIds := make([]string, len(fileMap))
	i := 0
	for k, v := range fileMap {
		blobIds[i] = k
		paths[i] = v
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
			FileMap: map[string]string{
				"123": "a/path/to/file",
				"456": "another/path/to/other/file",
				"789": "yet/again/a/path",
				"234": "i/am/getting/bored"}},
			[]string{"123", "789"}},
		{clientcommit.Commit{
			Message: "empty commit",
			Creator: "user",
			FileMap: map[string]string{}},
			[]string{}},
	}
	for _, testcase := range testcases {
		commFake := &netiofakes.FakeCommunicator{}
		fileFake := &fileiofakes.FakeLocalcopy{}
		paths, blobIds := filesFromMap(testcase.expectedCommit.FileMap)

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
			if actual, _ := fileFake.SendBlobArgsForCall(i); actual != testcase.expectedCommit.FileMap[blobId] {
				t.Errorf("Send string called with wrong args when sending blob id: Expected: %s, actual: %s", testcase.expectedCommit.FileMap[blobId], actual)
			}
		}

	}
}
