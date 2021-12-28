package clientUI

import (
	"bufio"
	"fmt"
	"net"
	"os"

	"github.com/NikolayDPaev/CentralisedVersionControl/client/commands"
)

func Init() error {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter username:")
	scanner.Scan()
	username := scanner.Text()

	fmt.Println("Enter remote address:")
	scanner.Scan()
	address := scanner.Text()

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading user input: %w", scanner.Err())
	}

	metafileData := commands.MetafileData{Username: username, Address: address}
	if err := metafileData.Save(); err != nil {
		return fmt.Errorf("error initializing: %w", scanner.Err())
	}
	return nil
}

func CommitList() error {
	metafile, err := commands.ReadMetafileData()
	if err != nil {
		return err
	}

	c, err := net.Dial("tcp", metafile.Address)
	if err != nil {
		return fmt.Errorf("error connecting to server: %w", err)
	}
	defer c.Close()

	commitList, err := commands.GetCommitList(c, c)
	if err != nil {
		return fmt.Errorf("cannot execute commit list operation: %w", err)
	}

	fmt.Println("Commits on server:")
	for _, str := range commitList {
		fmt.Println(str)
	}

	return nil
}

func DownloadCommit(commitId string) error {
	metafile, err := commands.ReadMetafileData()
	if err != nil {
		return err
	}

	c, err := net.Dial("tcp", metafile.Address)
	if err != nil {
		return fmt.Errorf("error connecting to server: %w", err)
	}
	defer c.Close()

	message, err := commands.DownloadCommit(commitId, c, c)
	if err != nil {
		return fmt.Errorf("cannot execute download commit operation: %w", err)
	}
	fmt.Println(message)
	return nil
}

func UploadCommit() error {
	metafile, err := commands.ReadMetafileData()
	if err != nil {
		return err
	}

	c, err := net.Dial("tcp", metafile.Address)
	if err != nil {
		return fmt.Errorf("error connecting to server: %w", err)
	}
	defer c.Close()

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter commit message:")
	scanner.Scan()
	message := scanner.Text()

	err = commands.UploadCommit(message, metafile.Username, c, c)
	if err != nil {
		return fmt.Errorf("error uploading commit: %w", err)
	}
	fmt.Println("Commit uploaded successfuly!")
	return nil
}
