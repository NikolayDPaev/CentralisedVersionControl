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
		return fmt.Errorf("Error reading user input: %w", scanner.Err())
	}

	metafileData := commands.MetafileData{Username: username, Address: address}
	if err := metafileData.Save(); err != nil {
		return fmt.Errorf("Error initializing: %w", scanner.Err())
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

	commitList, err := commands.GetCommitList(c)
	if err != nil {
		return fmt.Errorf("cannot execute commit list operation: %w", err)
	}

	fmt.Println("Commits on server:")
	for _, str := range commitList {
		fmt.Println(str)
	}

	return nil
}
