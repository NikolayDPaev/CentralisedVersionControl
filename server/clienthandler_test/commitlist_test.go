package clienthandler_test

import (
	"reflect"
	"testing"

	"github.com/NikolayDPaev/CentralisedVersionControl/server/clienthandler"
)

type commitListCommunicatorMock struct {
	communicatorDefMock
	testStrings []string
	successful  bool
}

type commitListStorageMock struct {
	storageDefMock
	testStrings []string
}

func (s *commitListStorageMock) CommitList() []string {
	return []string{"some", "strings", "representing", "commit"}
}

func (c *commitListCommunicatorMock) SendStringSlice(slice []string) error {
	if reflect.DeepEqual(slice, []string{"some", "strings", "representing", "commit"}) {
		c.successful = true
	}
	return nil
}

func TestCommitListHandle(t *testing.T) {
	tests := [][]string{
		{},
		{"some", "strings", "representing", "commit"},
	}
	for _, test := range tests {
		netMock := &commitListCommunicatorMock{}
		fileMock := &commitListStorageMock{}
		netMock.testStrings = test
		fileMock.testStrings = test

		commitList := clienthandler.NewCommitList(netMock, fileMock)
		commitList.Handle()

		if !netMock.successful {
			t.FailNow()
		}
	}
}
