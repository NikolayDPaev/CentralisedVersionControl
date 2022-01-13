package clienthandler_test

import (
	"reflect"
	"testing"

	"github.com/NikolayDPaev/CentralisedVersionControl/server/clienthandler"
)

type communicatorMock struct {
	communicatorDefMock
	successful bool
}

type storageMock struct {
	storageDefMock
}

func (s *storageMock) CommitList() []string {
	return []string{"some", "strings", "representing", "commit"}
}

func (c *communicatorMock) SendStringSlice(slice []string) error {
	if reflect.DeepEqual(slice, []string{"some", "strings", "representing", "commit"}) {
		c.successful = true
	}
	return nil
}

func TestHandle(t *testing.T) {
	netMock := &communicatorMock{}
	fileMock := &storageMock{}

	commitList := clienthandler.NewCommitList(netMock, fileMock)
	commitList.Handle()

	if !netMock.successful {
		t.FailNow()
	}
}
