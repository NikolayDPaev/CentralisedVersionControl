// Code generated by counterfeiter. DO NOT EDIT.
package fileiofakes

import (
	"sync"

	"github.com/NikolayDPaev/CentralisedVersionControl/client/fileio"
	"github.com/NikolayDPaev/CentralisedVersionControl/netio"
)

type FakeLocalcopy struct {
	CleanOtherFilesStub        func(map[string]struct{}) error
	cleanOtherFilesMutex       sync.RWMutex
	cleanOtherFilesArgsForCall []struct {
		arg1 map[string]struct{}
	}
	cleanOtherFilesReturns struct {
		result1 error
	}
	cleanOtherFilesReturnsOnCall map[int]struct {
		result1 error
	}
	FileSizeStub        func(string) (int64, error)
	fileSizeMutex       sync.RWMutex
	fileSizeArgsForCall []struct {
		arg1 string
	}
	fileSizeReturns struct {
		result1 int64
		result2 error
	}
	fileSizeReturnsOnCall map[int]struct {
		result1 int64
		result2 error
	}
	FileWithHashExistsStub        func(string, string) (bool, error)
	fileWithHashExistsMutex       sync.RWMutex
	fileWithHashExistsArgsForCall []struct {
		arg1 string
		arg2 string
	}
	fileWithHashExistsReturns struct {
		result1 bool
		result2 error
	}
	fileWithHashExistsReturnsOnCall map[int]struct {
		result1 bool
		result2 error
	}
	GetHashOfFileStub        func(string) (string, error)
	getHashOfFileMutex       sync.RWMutex
	getHashOfFileArgsForCall []struct {
		arg1 string
	}
	getHashOfFileReturns struct {
		result1 string
		result2 error
	}
	getHashOfFileReturnsOnCall map[int]struct {
		result1 string
		result2 error
	}
	GetPathsOfAllFilesStub        func() ([]string, error)
	getPathsOfAllFilesMutex       sync.RWMutex
	getPathsOfAllFilesArgsForCall []struct {
	}
	getPathsOfAllFilesReturns struct {
		result1 []string
		result2 error
	}
	getPathsOfAllFilesReturnsOnCall map[int]struct {
		result1 []string
		result2 error
	}
	ReceiveBlobStub        func(string, netio.Communicator) error
	receiveBlobMutex       sync.RWMutex
	receiveBlobArgsForCall []struct {
		arg1 string
		arg2 netio.Communicator
	}
	receiveBlobReturns struct {
		result1 error
	}
	receiveBlobReturnsOnCall map[int]struct {
		result1 error
	}
	SendBlobStub        func(string, netio.Communicator) error
	sendBlobMutex       sync.RWMutex
	sendBlobArgsForCall []struct {
		arg1 string
		arg2 netio.Communicator
	}
	sendBlobReturns struct {
		result1 error
	}
	sendBlobReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeLocalcopy) CleanOtherFiles(arg1 map[string]struct{}) error {
	fake.cleanOtherFilesMutex.Lock()
	ret, specificReturn := fake.cleanOtherFilesReturnsOnCall[len(fake.cleanOtherFilesArgsForCall)]
	fake.cleanOtherFilesArgsForCall = append(fake.cleanOtherFilesArgsForCall, struct {
		arg1 map[string]struct{}
	}{arg1})
	stub := fake.CleanOtherFilesStub
	fakeReturns := fake.cleanOtherFilesReturns
	fake.recordInvocation("CleanOtherFiles", []interface{}{arg1})
	fake.cleanOtherFilesMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeLocalcopy) CleanOtherFilesCallCount() int {
	fake.cleanOtherFilesMutex.RLock()
	defer fake.cleanOtherFilesMutex.RUnlock()
	return len(fake.cleanOtherFilesArgsForCall)
}

func (fake *FakeLocalcopy) CleanOtherFilesCalls(stub func(map[string]struct{}) error) {
	fake.cleanOtherFilesMutex.Lock()
	defer fake.cleanOtherFilesMutex.Unlock()
	fake.CleanOtherFilesStub = stub
}

func (fake *FakeLocalcopy) CleanOtherFilesArgsForCall(i int) map[string]struct{} {
	fake.cleanOtherFilesMutex.RLock()
	defer fake.cleanOtherFilesMutex.RUnlock()
	argsForCall := fake.cleanOtherFilesArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeLocalcopy) CleanOtherFilesReturns(result1 error) {
	fake.cleanOtherFilesMutex.Lock()
	defer fake.cleanOtherFilesMutex.Unlock()
	fake.CleanOtherFilesStub = nil
	fake.cleanOtherFilesReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeLocalcopy) CleanOtherFilesReturnsOnCall(i int, result1 error) {
	fake.cleanOtherFilesMutex.Lock()
	defer fake.cleanOtherFilesMutex.Unlock()
	fake.CleanOtherFilesStub = nil
	if fake.cleanOtherFilesReturnsOnCall == nil {
		fake.cleanOtherFilesReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.cleanOtherFilesReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeLocalcopy) FileSize(arg1 string) (int64, error) {
	fake.fileSizeMutex.Lock()
	ret, specificReturn := fake.fileSizeReturnsOnCall[len(fake.fileSizeArgsForCall)]
	fake.fileSizeArgsForCall = append(fake.fileSizeArgsForCall, struct {
		arg1 string
	}{arg1})
	stub := fake.FileSizeStub
	fakeReturns := fake.fileSizeReturns
	fake.recordInvocation("FileSize", []interface{}{arg1})
	fake.fileSizeMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeLocalcopy) FileSizeCallCount() int {
	fake.fileSizeMutex.RLock()
	defer fake.fileSizeMutex.RUnlock()
	return len(fake.fileSizeArgsForCall)
}

func (fake *FakeLocalcopy) FileSizeCalls(stub func(string) (int64, error)) {
	fake.fileSizeMutex.Lock()
	defer fake.fileSizeMutex.Unlock()
	fake.FileSizeStub = stub
}

func (fake *FakeLocalcopy) FileSizeArgsForCall(i int) string {
	fake.fileSizeMutex.RLock()
	defer fake.fileSizeMutex.RUnlock()
	argsForCall := fake.fileSizeArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeLocalcopy) FileSizeReturns(result1 int64, result2 error) {
	fake.fileSizeMutex.Lock()
	defer fake.fileSizeMutex.Unlock()
	fake.FileSizeStub = nil
	fake.fileSizeReturns = struct {
		result1 int64
		result2 error
	}{result1, result2}
}

func (fake *FakeLocalcopy) FileSizeReturnsOnCall(i int, result1 int64, result2 error) {
	fake.fileSizeMutex.Lock()
	defer fake.fileSizeMutex.Unlock()
	fake.FileSizeStub = nil
	if fake.fileSizeReturnsOnCall == nil {
		fake.fileSizeReturnsOnCall = make(map[int]struct {
			result1 int64
			result2 error
		})
	}
	fake.fileSizeReturnsOnCall[i] = struct {
		result1 int64
		result2 error
	}{result1, result2}
}

func (fake *FakeLocalcopy) FileWithHashExists(arg1 string, arg2 string) (bool, error) {
	fake.fileWithHashExistsMutex.Lock()
	ret, specificReturn := fake.fileWithHashExistsReturnsOnCall[len(fake.fileWithHashExistsArgsForCall)]
	fake.fileWithHashExistsArgsForCall = append(fake.fileWithHashExistsArgsForCall, struct {
		arg1 string
		arg2 string
	}{arg1, arg2})
	stub := fake.FileWithHashExistsStub
	fakeReturns := fake.fileWithHashExistsReturns
	fake.recordInvocation("FileWithHashExists", []interface{}{arg1, arg2})
	fake.fileWithHashExistsMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeLocalcopy) FileWithHashExistsCallCount() int {
	fake.fileWithHashExistsMutex.RLock()
	defer fake.fileWithHashExistsMutex.RUnlock()
	return len(fake.fileWithHashExistsArgsForCall)
}

func (fake *FakeLocalcopy) FileWithHashExistsCalls(stub func(string, string) (bool, error)) {
	fake.fileWithHashExistsMutex.Lock()
	defer fake.fileWithHashExistsMutex.Unlock()
	fake.FileWithHashExistsStub = stub
}

func (fake *FakeLocalcopy) FileWithHashExistsArgsForCall(i int) (string, string) {
	fake.fileWithHashExistsMutex.RLock()
	defer fake.fileWithHashExistsMutex.RUnlock()
	argsForCall := fake.fileWithHashExistsArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeLocalcopy) FileWithHashExistsReturns(result1 bool, result2 error) {
	fake.fileWithHashExistsMutex.Lock()
	defer fake.fileWithHashExistsMutex.Unlock()
	fake.FileWithHashExistsStub = nil
	fake.fileWithHashExistsReturns = struct {
		result1 bool
		result2 error
	}{result1, result2}
}

func (fake *FakeLocalcopy) FileWithHashExistsReturnsOnCall(i int, result1 bool, result2 error) {
	fake.fileWithHashExistsMutex.Lock()
	defer fake.fileWithHashExistsMutex.Unlock()
	fake.FileWithHashExistsStub = nil
	if fake.fileWithHashExistsReturnsOnCall == nil {
		fake.fileWithHashExistsReturnsOnCall = make(map[int]struct {
			result1 bool
			result2 error
		})
	}
	fake.fileWithHashExistsReturnsOnCall[i] = struct {
		result1 bool
		result2 error
	}{result1, result2}
}

func (fake *FakeLocalcopy) GetHashOfFile(arg1 string) (string, error) {
	fake.getHashOfFileMutex.Lock()
	ret, specificReturn := fake.getHashOfFileReturnsOnCall[len(fake.getHashOfFileArgsForCall)]
	fake.getHashOfFileArgsForCall = append(fake.getHashOfFileArgsForCall, struct {
		arg1 string
	}{arg1})
	stub := fake.GetHashOfFileStub
	fakeReturns := fake.getHashOfFileReturns
	fake.recordInvocation("GetHashOfFile", []interface{}{arg1})
	fake.getHashOfFileMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeLocalcopy) GetHashOfFileCallCount() int {
	fake.getHashOfFileMutex.RLock()
	defer fake.getHashOfFileMutex.RUnlock()
	return len(fake.getHashOfFileArgsForCall)
}

func (fake *FakeLocalcopy) GetHashOfFileCalls(stub func(string) (string, error)) {
	fake.getHashOfFileMutex.Lock()
	defer fake.getHashOfFileMutex.Unlock()
	fake.GetHashOfFileStub = stub
}

func (fake *FakeLocalcopy) GetHashOfFileArgsForCall(i int) string {
	fake.getHashOfFileMutex.RLock()
	defer fake.getHashOfFileMutex.RUnlock()
	argsForCall := fake.getHashOfFileArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeLocalcopy) GetHashOfFileReturns(result1 string, result2 error) {
	fake.getHashOfFileMutex.Lock()
	defer fake.getHashOfFileMutex.Unlock()
	fake.GetHashOfFileStub = nil
	fake.getHashOfFileReturns = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeLocalcopy) GetHashOfFileReturnsOnCall(i int, result1 string, result2 error) {
	fake.getHashOfFileMutex.Lock()
	defer fake.getHashOfFileMutex.Unlock()
	fake.GetHashOfFileStub = nil
	if fake.getHashOfFileReturnsOnCall == nil {
		fake.getHashOfFileReturnsOnCall = make(map[int]struct {
			result1 string
			result2 error
		})
	}
	fake.getHashOfFileReturnsOnCall[i] = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeLocalcopy) GetPathsOfAllFiles() ([]string, error) {
	fake.getPathsOfAllFilesMutex.Lock()
	ret, specificReturn := fake.getPathsOfAllFilesReturnsOnCall[len(fake.getPathsOfAllFilesArgsForCall)]
	fake.getPathsOfAllFilesArgsForCall = append(fake.getPathsOfAllFilesArgsForCall, struct {
	}{})
	stub := fake.GetPathsOfAllFilesStub
	fakeReturns := fake.getPathsOfAllFilesReturns
	fake.recordInvocation("GetPathsOfAllFiles", []interface{}{})
	fake.getPathsOfAllFilesMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeLocalcopy) GetPathsOfAllFilesCallCount() int {
	fake.getPathsOfAllFilesMutex.RLock()
	defer fake.getPathsOfAllFilesMutex.RUnlock()
	return len(fake.getPathsOfAllFilesArgsForCall)
}

func (fake *FakeLocalcopy) GetPathsOfAllFilesCalls(stub func() ([]string, error)) {
	fake.getPathsOfAllFilesMutex.Lock()
	defer fake.getPathsOfAllFilesMutex.Unlock()
	fake.GetPathsOfAllFilesStub = stub
}

func (fake *FakeLocalcopy) GetPathsOfAllFilesReturns(result1 []string, result2 error) {
	fake.getPathsOfAllFilesMutex.Lock()
	defer fake.getPathsOfAllFilesMutex.Unlock()
	fake.GetPathsOfAllFilesStub = nil
	fake.getPathsOfAllFilesReturns = struct {
		result1 []string
		result2 error
	}{result1, result2}
}

func (fake *FakeLocalcopy) GetPathsOfAllFilesReturnsOnCall(i int, result1 []string, result2 error) {
	fake.getPathsOfAllFilesMutex.Lock()
	defer fake.getPathsOfAllFilesMutex.Unlock()
	fake.GetPathsOfAllFilesStub = nil
	if fake.getPathsOfAllFilesReturnsOnCall == nil {
		fake.getPathsOfAllFilesReturnsOnCall = make(map[int]struct {
			result1 []string
			result2 error
		})
	}
	fake.getPathsOfAllFilesReturnsOnCall[i] = struct {
		result1 []string
		result2 error
	}{result1, result2}
}

func (fake *FakeLocalcopy) ReceiveBlob(arg1 string, arg2 netio.Communicator) error {
	fake.receiveBlobMutex.Lock()
	ret, specificReturn := fake.receiveBlobReturnsOnCall[len(fake.receiveBlobArgsForCall)]
	fake.receiveBlobArgsForCall = append(fake.receiveBlobArgsForCall, struct {
		arg1 string
		arg2 netio.Communicator
	}{arg1, arg2})
	stub := fake.ReceiveBlobStub
	fakeReturns := fake.receiveBlobReturns
	fake.recordInvocation("ReceiveBlob", []interface{}{arg1, arg2})
	fake.receiveBlobMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeLocalcopy) ReceiveBlobCallCount() int {
	fake.receiveBlobMutex.RLock()
	defer fake.receiveBlobMutex.RUnlock()
	return len(fake.receiveBlobArgsForCall)
}

func (fake *FakeLocalcopy) ReceiveBlobCalls(stub func(string, netio.Communicator) error) {
	fake.receiveBlobMutex.Lock()
	defer fake.receiveBlobMutex.Unlock()
	fake.ReceiveBlobStub = stub
}

func (fake *FakeLocalcopy) ReceiveBlobArgsForCall(i int) (string, netio.Communicator) {
	fake.receiveBlobMutex.RLock()
	defer fake.receiveBlobMutex.RUnlock()
	argsForCall := fake.receiveBlobArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeLocalcopy) ReceiveBlobReturns(result1 error) {
	fake.receiveBlobMutex.Lock()
	defer fake.receiveBlobMutex.Unlock()
	fake.ReceiveBlobStub = nil
	fake.receiveBlobReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeLocalcopy) ReceiveBlobReturnsOnCall(i int, result1 error) {
	fake.receiveBlobMutex.Lock()
	defer fake.receiveBlobMutex.Unlock()
	fake.ReceiveBlobStub = nil
	if fake.receiveBlobReturnsOnCall == nil {
		fake.receiveBlobReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.receiveBlobReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeLocalcopy) SendBlob(arg1 string, arg2 netio.Communicator) error {
	fake.sendBlobMutex.Lock()
	ret, specificReturn := fake.sendBlobReturnsOnCall[len(fake.sendBlobArgsForCall)]
	fake.sendBlobArgsForCall = append(fake.sendBlobArgsForCall, struct {
		arg1 string
		arg2 netio.Communicator
	}{arg1, arg2})
	stub := fake.SendBlobStub
	fakeReturns := fake.sendBlobReturns
	fake.recordInvocation("SendBlob", []interface{}{arg1, arg2})
	fake.sendBlobMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeLocalcopy) SendBlobCallCount() int {
	fake.sendBlobMutex.RLock()
	defer fake.sendBlobMutex.RUnlock()
	return len(fake.sendBlobArgsForCall)
}

func (fake *FakeLocalcopy) SendBlobCalls(stub func(string, netio.Communicator) error) {
	fake.sendBlobMutex.Lock()
	defer fake.sendBlobMutex.Unlock()
	fake.SendBlobStub = stub
}

func (fake *FakeLocalcopy) SendBlobArgsForCall(i int) (string, netio.Communicator) {
	fake.sendBlobMutex.RLock()
	defer fake.sendBlobMutex.RUnlock()
	argsForCall := fake.sendBlobArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeLocalcopy) SendBlobReturns(result1 error) {
	fake.sendBlobMutex.Lock()
	defer fake.sendBlobMutex.Unlock()
	fake.SendBlobStub = nil
	fake.sendBlobReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeLocalcopy) SendBlobReturnsOnCall(i int, result1 error) {
	fake.sendBlobMutex.Lock()
	defer fake.sendBlobMutex.Unlock()
	fake.SendBlobStub = nil
	if fake.sendBlobReturnsOnCall == nil {
		fake.sendBlobReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.sendBlobReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeLocalcopy) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.cleanOtherFilesMutex.RLock()
	defer fake.cleanOtherFilesMutex.RUnlock()
	fake.fileSizeMutex.RLock()
	defer fake.fileSizeMutex.RUnlock()
	fake.fileWithHashExistsMutex.RLock()
	defer fake.fileWithHashExistsMutex.RUnlock()
	fake.getHashOfFileMutex.RLock()
	defer fake.getHashOfFileMutex.RUnlock()
	fake.getPathsOfAllFilesMutex.RLock()
	defer fake.getPathsOfAllFilesMutex.RUnlock()
	fake.receiveBlobMutex.RLock()
	defer fake.receiveBlobMutex.RUnlock()
	fake.sendBlobMutex.RLock()
	defer fake.sendBlobMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeLocalcopy) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ fileio.Localcopy = new(FakeLocalcopy)
