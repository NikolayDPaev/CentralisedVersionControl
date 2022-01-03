package commands

import (
	"fmt"
	"io"

	"github.com/NikolayDPaev/CentralisedVersionControl/client/commit"
	"github.com/NikolayDPaev/CentralisedVersionControl/client/fileIO"
	"github.com/NikolayDPaev/CentralisedVersionControl/client/netIO"
)

func sendCommit(commit *commit.Commit, writer io.Writer) error {
	err := netIO.SendString(commit.Md5Hash(), writer)
	if err != nil {
		return fmt.Errorf("error sending commit id:\n%w", err)
	}

	err = commit.Send(writer)
	if err != nil {
		return fmt.Errorf("error sending commit:\n%w", err)
	}
	return nil
}

func sendCompressedBlob(filePath string, writer io.Writer) error {
	tmpFile, err := fileIO.CompressToTempFile(filePath)
	if err != nil {
		return fmt.Errorf("error compressing file %s:\n%w", filePath, err)
	}
	defer tmpFile.Close()

	stat, err := tmpFile.Stat()
	if err != nil {
		return fmt.Errorf("error getting blobTmp size:\n%w", err)
	}

	err = netIO.SendFileData(tmpFile, stat.Size(), writer)
	if err != nil {
		return fmt.Errorf("error sending blob:\n%w", err)
	}
	return nil
}

func sendBlob(blobId string, commit *commit.Commit, writer io.Writer) error {
	if err := netIO.SendString(blobId, writer); err != nil {
		return fmt.Errorf("error sending blobId %s:\n%w", blobId, err)
	}

	path, err := commit.GetBlobPath(blobId)
	if err != nil {
		return err
	}

	if err := sendCompressedBlob(path, writer); err != nil {
		return fmt.Errorf("error sending file %s:\n%w", path, err)
	}
	return nil
}

func UploadCommit(message, username string, reader io.Reader, writer io.Writer) error {
	commit, err := commit.CreateCommit(message, username)
	if err != nil {
		return fmt.Errorf("error creating commit:\n%w", err)
	}

	fmt.Printf("Sending commit %s\n", commit.Md5Hash())
	err = sendCommit(commit, writer)
	if err != nil {
		return err
	}

	fmt.Println("Receiving missing blob IDs from server")
	missingBlobIds, err := netIO.ReceiveStringSlice(reader)
	if err != nil {
		return fmt.Errorf("error receiving missing blobIds:\n%w", err)
	}

	fmt.Printf("Sending %d objects\n", len(missingBlobIds))
	for _, blobId := range missingBlobIds {
		err := sendBlob(blobId, commit, writer)
		if err != nil {
			return fmt.Errorf("error sending blob with id %s:\n%w", blobId, err)
		}
	}
	return nil
}
