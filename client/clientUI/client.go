package clientUI

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/NikolayDPaev/CentralisedVersionControl/client/commands"
	"github.com/NikolayDPaev/CentralisedVersionControl/client/netIO"
)

const (
	GET_COMMIT_LIST = 0
	DOWNLOAD_COMMIT = 1
	UPLOAD_COMMIT   = 2
	EMPTY_REQUEST   = 3
)

var errInvalidCommand = errors.New("invalid command")

func ReadArgs(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: csv <command>")
		return
	}

	var err error
	switch args[0] {
	case "init":
		err = initClient()
	case "list":
		err = commitList()
	case "push":
		err = uploadCommit()
	case "pull":
		err = downloadCommit(args)
	case "help":
		help()
	default:
		err = errInvalidCommand
	}

	if err != nil {
		if errors.Is(err, commands.ErrMissingMetafile) {
			fmt.Println("Cannot find .cvc file! Please run command csv init.")
		} else if errors.Is(err, errInvalidCommand) {
			fmt.Println("Incorrect command. For list of commands run \"csv help\".")
		} else {
			log.Println(err)
		}
	}
}

func help() {
	fmt.Println("Usage: csv <command>")
	fmt.Println("Commands:")
	fmt.Println("init - initialization of workplace in the current directory")
	fmt.Println("list - listing of the available commits")
	fmt.Println("pull <commitId> - downloading commit to the workplace")
	fmt.Println("push - commiting the current state of the workplace ")
	fmt.Println("help - prints this text")
}

func initClient() error {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter username:")
	scanner.Scan()
	username := scanner.Text()

	fmt.Println("Enter remote address:")
	scanner.Scan()
	address := scanner.Text()

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading user input:\n%w", scanner.Err())
	}

	metafileData := commands.MetafileData{Username: username, Address: address}
	if err := metafileData.Save(); err != nil {
		return fmt.Errorf("error initializing:\n%w", scanner.Err())
	}
	return nil
}

func commitList() error {
	metafile, err := commands.ReadMetafileData()
	if err != nil {
		return err
	}

	c, err := net.Dial("tcp", metafile.Address)
	if err != nil {
		return fmt.Errorf("error connecting to server:\n%w", err)
	}
	defer c.Close()

	if err := netIO.SendVarInt(GET_COMMIT_LIST, c); err != nil {
		return fmt.Errorf("cannot send opcode:\n%w", err)
	}

	commitList, err := commands.GetCommitList(c, c)
	if err != nil {
		return fmt.Errorf("cannot execute commit list operation:\n%w", err)
	}

	fmt.Println("Commits on server:")
	for _, str := range commitList {
		fmt.Println(str)
	}

	return nil
}

func downloadCommit(args []string) error {
	if len(args) != 2 {
		return errInvalidCommand
	}
	metafile, err := commands.ReadMetafileData()
	if err != nil {
		return err
	}

	c, err := net.Dial("tcp", metafile.Address)
	if err != nil {
		return fmt.Errorf("error connecting to server:\n%w", err)
	}
	defer c.Close()

	if err := netIO.SendVarInt(DOWNLOAD_COMMIT, c); err != nil {
		return fmt.Errorf("cannot send opcode:\n%w", err)
	}

	message, err := commands.DownloadCommit(args[1], c, c)
	if err != nil {
		return fmt.Errorf("cannot execute download commit operation:\n%w", err)
	}
	fmt.Println(message)
	return nil
}

func uploadCommit() error {
	metafile, err := commands.ReadMetafileData()
	if err != nil {
		return err
	}

	c, err := net.Dial("tcp", metafile.Address)
	if err != nil {
		return fmt.Errorf("error connecting to server:\n%w", err)
	}
	defer c.Close()

	if err := netIO.SendVarInt(UPLOAD_COMMIT, c); err != nil {
		return fmt.Errorf("cannot send opcode:\n%w", err)
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter commit message:")
	scanner.Scan()
	message := scanner.Text()

	err = commands.UploadCommit(message, metafile.Username, c, c)
	if err != nil {
		return fmt.Errorf("error uploading commit:\n%w", err)
	}
	fmt.Println("Commit uploaded successfuly!")
	return nil
}
