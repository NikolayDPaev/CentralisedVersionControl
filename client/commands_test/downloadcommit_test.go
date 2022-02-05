package commands_test

import (
	"testing"

	"github.com/NikolayDPaev/CentralisedVersionControl/client/clientcommit"
	"github.com/NikolayDPaev/CentralisedVersionControl/client/commands"
	"github.com/NikolayDPaev/CentralisedVersionControl/client/fileio/fileiofakes"
	"github.com/NikolayDPaev/CentralisedVersionControl/netio/netiofakes"
)

type downloadTestCase struct {
	message           string
	creator           string
	tree              string
	paths             []string
	missingBlobsPaths []string
	missingBlobs      []string
}

const (
	OPCODE = 1
	OK     = 1
	ERROR  = 0
)

func equalSets(s1, s2 []string) bool {
	for _, a := range s1 {
		member := false
		for _, b := range s2 {
			if a == b {
				member = true
			}
		}
		if !member {
			return false
		}
	}
	return true
}

func TestDownloadcommit(t *testing.T) {
	testCases := []downloadTestCase{
		{"random message",
			"creator",
			"123 /some/path/to/blob\n456 /different/path\n789 /another/path\n101 /and/another",
			[]string{"/some/path/to/blob", "/different/path", "/and/another", "/another/path"},
			[]string{"/another/path", "/some/path/to/blob"},
			[]string{"789", "123"}},
		{"empty commit",
			"user",
			"",
			[]string{},
			[]string{},
			[]string{}},
	}
	for _, testCase := range testCases {
		filesMap, _ := clientcommit.GetMap(testCase.tree)
		testCommit := &clientcommit.Commit{Message: testCase.message, Creator: testCase.creator, FileMap: filesMap}
		testCommitId := testCommit.Md5Hash()

		commFake := &netiofakes.FakeCommunicator{}
		fileFake := &fileiofakes.FakeLocalcopy{}

		commFake.RecvVarIntReturnsOnCall(0, OK, nil)
		commFake.RecvStringReturnsOnCall(0, testCommitId, nil)
		commFake.RecvStringReturnsOnCall(1, testCommit.Message, nil)
		commFake.RecvStringReturnsOnCall(2, testCommit.Creator, nil)
		commFake.RecvStringReturnsOnCall(3, testCase.tree, nil)
		fileFake.FileWithHashExistsStub = func(s1, s2 string) (bool, error) {
			for _, id := range testCase.missingBlobs {
				if s2 == id {
					return false, nil
				}
			}
			return true, nil
		}

		for i, blob := range testCase.missingBlobs {
			commFake.RecvStringReturnsOnCall(4+i, blob, nil)
		}

		downloadCommit := commands.NewDownload(commFake, fileFake, OPCODE, OK)
		if err := downloadCommit.DownloadCommit(testCommitId); err != nil {
			t.Errorf("Error catched: %s", err)
		}

		if actual := commFake.SendVarIntArgsForCall(0); actual != OPCODE {
			t.Errorf("Send var int called with wrong arg, when sending opcode: Expected %d, actual %d", OPCODE, actual)
		}

		if actual := commFake.SendStringArgsForCall(0); actual != testCommitId {
			t.Errorf("Send string called with wrong arg, when sending commitId: Expected %s, actual %s", testCommitId, actual)
		}

		if len(testCase.missingBlobs) > 0 {
			if actual := commFake.SendStringSliceArgsForCall(0); !equalSets(actual, testCase.missingBlobs) {
				t.Errorf("Send string slice called with wrong arg when sending missing blobs: Expected %s, actual %s", testCase.missingBlobs, actual)
			}
		}

		actualPathSet := fileFake.CleanOtherFilesArgsForCall(0)
		for _, str := range testCase.paths {
			if _, ok := actualPathSet[str]; !ok {
				t.Errorf("Clean other files called with wrong args: Expected %s, actual %s", testCase.paths, actualPathSet)
			}
		}

		for i, missingPaths := range testCase.missingBlobsPaths {
			if actual, _ := fileFake.ReceiveBlobArgsForCall(i); actual != missingPaths {
				t.Errorf("Receive blob called with wrong arg when receiving missing blob: Expected %s, actual %s", missingPaths, actual)
			}
		}
	}
}
