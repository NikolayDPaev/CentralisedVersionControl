package clienthandler_test

import (
	"reflect"
	"testing"

	"github.com/NikolayDPaev/CentralisedVersionControl/server/clienthandler"
	"github.com/NikolayDPaev/CentralisedVersionControl/server/commit"
	"github.com/NikolayDPaev/CentralisedVersionControl/server/netIO"
)

type receivedData struct {
	stringsToReceive   []string
	wantedMissingBlobs []string
}

type receiveCommitCommunicatorMock struct {
	communicatorDefMock

	stringsToReceive   []string
	receivedCounter    int
	wantedMissingBlobs []string
	t                  *testing.T
}

func (c *receiveCommitCommunicatorMock) ReceiveString() (string, error) {
	str := c.stringsToReceive[c.receivedCounter]
	c.receivedCounter++
	return str, nil
}

func (c *receiveCommitCommunicatorMock) SendStringSlice(slice []string) error {
	if (len(slice) != 0 || len(c.wantedMissingBlobs) != 0) && !reflect.DeepEqual(slice, c.wantedMissingBlobs) {
		c.t.Errorf("Send string failed: expected: %s actual: %s", c.wantedMissingBlobs, slice)
	}
	return nil
}

type receiveCommitStorageMock struct {
	storageDefMock

	stringsToReceive   []string
	wantedMissingBlobs []string
	t                  *testing.T
}

func (s *receiveCommitStorageMock) BlobExists(blobId string) (bool, error) {
	for _, v := range s.wantedMissingBlobs {
		if blobId == v {
			return false, nil
		}
	}
	return true, nil
}

func (s *receiveCommitStorageMock) SaveBlob(blobId string, comm netIO.Communicator) error {
	for _, v := range s.wantedMissingBlobs {
		if blobId == v {
			return nil
		}
	}
	s.t.Errorf("Send blob failed: save blob was not called with any of: %s, instead: %s", s.wantedMissingBlobs, blobId)
	return nil
}

func (s *receiveCommitStorageMock) SaveCommit(receivedCommit *commit.Commit) error {
	wantedCommit := commit.NewCommit(s.stringsToReceive[0],
		s.stringsToReceive[1],
		s.stringsToReceive[2],
		s.stringsToReceive[3])
	if !reflect.DeepEqual(wantedCommit, receivedCommit) {
		s.t.Errorf("Send commit failed: commits does not match expected: %s, instead: %s", wantedCommit, receivedCommit)
	}
	return nil
}

func TestReceiveCommitHandle(t *testing.T) {
	var testCases = []receivedData{
		{[]string{"12345", "message", "creator",
			"123 /some/path/to/blob\n456 /different/path\n789 /another/path\n101 /and/another", // commit data
			"789", "456"}, // sended blobs
			[]string{"456", "789"}}, // missing blobs

		{[]string{"12345", "empty commit", "creator",
			""}, // empty commit data
			[]string{}},
	}
	for _, testCase := range testCases {
		netMock := &receiveCommitCommunicatorMock{communicatorDefMock{}, testCase.stringsToReceive, 0, testCase.wantedMissingBlobs, t}
		fileMock := &receiveCommitStorageMock{storageDefMock{}, testCase.stringsToReceive, testCase.wantedMissingBlobs, t}

		fileMock.stringsToReceive = testCase.stringsToReceive
		fileMock.wantedMissingBlobs = testCase.wantedMissingBlobs

		commitList := clienthandler.NewReceiveCommit(netMock, fileMock)
		commitList.Handle()
	}
}
