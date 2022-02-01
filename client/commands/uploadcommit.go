package commands

import (
	"fmt"

	"github.com/NikolayDPaev/CentralisedVersionControl/client/clientcommit"
	"github.com/NikolayDPaev/CentralisedVersionControl/client/fileIO"
	"github.com/NikolayDPaev/CentralisedVersionControl/netIO"
)

type Upload struct {
	comm   netIO.Communicator
	opcode int
}

func NewUpload(comm netIO.Communicator, opcode int) *Upload {
	return &Upload{comm, opcode}
}

func (u *Upload) sendCommit(commit *clientcommit.Commit) error {
	err := u.comm.SendString(commit.Md5Hash())
	if err != nil {
		return fmt.Errorf("error sending commit id:\n%w", err)
	}

	err = commit.Send(u.comm)
	if err != nil {
		return fmt.Errorf("error sending commit:\n%w", err)
	}
	return nil
}

func (u *Upload) sendCompressedBlob(filePath string) error {
	tmpFile, err := fileIO.CompressToTempFile(filePath)
	if err != nil {
		return fmt.Errorf("error compressing file %s:\n%w", filePath, err)
	}
	defer tmpFile.Close()

	stat, err := tmpFile.Stat()
	if err != nil {
		return fmt.Errorf("error getting blobTmp size:\n%w", err)
	}

	err = u.comm.SendFileData(tmpFile, stat.Size())
	if err != nil {
		return fmt.Errorf("error sending blob:\n%w", err)
	}
	return nil
}

func (u *Upload) sendBlob(blobId string, commit *clientcommit.Commit) error {
	if err := u.comm.SendString(blobId); err != nil {
		return fmt.Errorf("error sending blobId %s:\n%w", blobId, err)
	}

	path, err := commit.GetBlobPath(blobId)
	if err != nil {
		return err
	}

	if err := u.sendCompressedBlob(path); err != nil {
		return fmt.Errorf("error sending file %s:\n%w", path, err)
	}
	return nil
}

func (u *Upload) UploadCommit(message, username string) error {
	if err := u.comm.SendVarInt(int64(u.opcode)); err != nil {
		return fmt.Errorf("cannot send opcode:\n%w", err)
	}

	commit, err := clientcommit.CreateCommit(message, username)
	if err != nil {
		return fmt.Errorf("error creating commit:\n%w", err)
	}

	fmt.Printf("Sending commit %s\n", commit.Md5Hash())
	err = u.sendCommit(commit)
	if err != nil {
		return err
	}

	fmt.Println("Receiving missing blob IDs from server")
	missingBlobIds, err := u.comm.ReceiveStringSlice()
	if err != nil {
		return fmt.Errorf("error receiving missing blobIds:\n%w", err)
	}

	fmt.Printf("Sending %d objects\n", len(missingBlobIds))
	for _, blobId := range missingBlobIds {
		err := u.sendBlob(blobId, commit)
		if err != nil {
			return fmt.Errorf("error sending blob with id %s:\n%w", blobId, err)
		}
	}
	return nil
}
