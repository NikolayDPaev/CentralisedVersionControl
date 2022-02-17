// Client provides the user interface of the client app
package client

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"path"

	"github.com/NikolayDPaev/CentralisedVersionControl/client/commands"
	"github.com/NikolayDPaev/CentralisedVersionControl/client/fileio"
	"github.com/NikolayDPaev/CentralisedVersionControl/client/metadata"
	"github.com/NikolayDPaev/CentralisedVersionControl/netio"
)

const (
	OK              = 0
	GET_COMMIT_LIST = 0
	DOWNLOAD_COMMIT = 1
	UPLOAD_COMMIT   = 2
	EMPTY_REQUEST   = 3
	CHUNK_SIZE      = 4096
)

var errInvalidCommand = errors.New("invalid command")
var errMetafileAlreadyExists = errors.New("metafile already exist")

const METAFILE_NAME = "./.cvc"

// ReadArgs is a function that receives the user args with which the app is called.
// It invokes the different commands based on the args provided.
// The function excepts args slice starting from the first argument
// (the zeroth - the name of the binary should be omitted - ReadArgs(os.Args[1:])).
func ReadArgs(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: cvc <command>")
		return
	}

	var err error
	switch args[0] {
	case "init":
		err = initClient()
	case "list":
		err = commitList()
	case "push":
		err = uploadCommit(args)
	case "pull":
		err = downloadCommit(args)
	case "help":
		help()
	default:
		err = errInvalidCommand
	}

	if err != nil {
		if errors.Is(err, metadata.ErrMissingMetafile) {
			fmt.Println("Cannot find .cvc file! Please run command cvc init.")
		} else if errors.Is(err, errMetafileAlreadyExists) {
			fmt.Println("Metafile already exists. Please remove it and then run cvc init.")
		} else if errors.Is(err, errInvalidCommand) {
			fmt.Println("Incorrect command. For list of commands run \"cvc help\".")
		} else {
			log.Println(err)
		}
	}
}

// help prints the help menu.
func help() {
	fmt.Println("Usage: cvc <command>")
	fmt.Println("Commands:")
	fmt.Println("init - initializes workplace in the current directory")
	fmt.Println("list - lists of the available commits")
	fmt.Println("pull <commitId> - downloads commit to the workplace")
	fmt.Println("push <message>- commits the current state of the workplace ")
	fmt.Println("help - prints this text")
}

// initClient initializes the client.
// Creates the metafile that stores username, remote address and ignored files.
// METAFILE and the client app binary (if its is present in the directory) are added as
// ignored files.
func initClient() error {
	if ok, err := fileio.FileExists(METAFILE_NAME); ok && err == nil {
		return errMetafileAlreadyExists
	}

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

	IgnoredFiles := make(map[string]struct{}, 2)
	IgnoredFiles[METAFILE_NAME] = struct{}{}

	if ok, err := fileio.FileExists(path.Base(os.Args[0])); ok && err == nil {
		IgnoredFiles[os.Args[0]] = struct{}{}
	}

	metafileData := &metadata.MetafileData{Username: username, Address: address, IgnoredFiles: IgnoredFiles}
	if err := metadata.Save(metafileData, METAFILE_NAME); err != nil {
		return fmt.Errorf("error initializing: %w", scanner.Err())
	}
	return nil
}

// prepCommunication reads metafile and connects to the server.
// Returns error if any of the two operation fails.
func prepCommunication() (*metadata.MetafileData, net.Conn, error) {
	metafile, err := metadata.ReadMetafileData(METAFILE_NAME)
	if err != nil {
		return nil, nil, err
	}

	c, err := net.Dial("tcp", metafile.Address)
	if err != nil {
		return nil, nil, fmt.Errorf("error connecting to server: %w", err)
	}
	return metafile, c, nil
}

// commitList attempts to make request for the commit list of the server
// by invoking the GetCommitlist method of the struct CommitList in commands package
// If successful - prints it.
func commitList() error {
	_, comm, err := prepCommunication()
	if err != nil {
		return err
	}
	defer comm.Close()

	commitList := commands.NewCommitList(netio.NewCommunication(CHUNK_SIZE, comm, comm), GET_COMMIT_LIST)
	slice, err := commitList.GetCommitList()
	if err != nil {
		return fmt.Errorf("cannot execute commit list operation: %w", err)
	}

	for _, str := range slice {
		fmt.Println(str)
	}

	return nil
}

// downloadCommit request a commit from the server by invoking
// the Download commit method of the struct Download in commands package
func downloadCommit(args []string) error {
	if len(args) != 2 {
		return errInvalidCommand
	}
	metafile, comm, err := prepCommunication()
	if err != nil {
		return err
	}
	defer comm.Close()

	download := commands.NewDownload(
		netio.NewCommunication(CHUNK_SIZE, comm, comm),
		fileio.NewLocalfiles(metafile.IgnoredFiles),
		DOWNLOAD_COMMIT,
		OK)

	if err := download.DownloadCommit(args[1]); err != nil {
		if errors.Is(err, commands.ErrInvalidCommitId) {
			fmt.Println("Invalid commit Id")
		} else {
			return fmt.Errorf("cannot execute download commit operation: %w", err)
		}
	}
	return nil
}

// uploadCommit attempts to add the current files as a commit and to send them to the server.
// Invokes the UploadCommit methods of Upload struct in commands package.
func uploadCommit(args []string) error {
	if len(args) != 2 {
		return errInvalidCommand
	}
	metafile, comm, err := prepCommunication()
	if err != nil {
		return err
	}
	defer comm.Close()

	upload := commands.NewUpload(
		netio.NewCommunication(CHUNK_SIZE, comm, comm),
		fileio.NewLocalfiles(metafile.IgnoredFiles),
		UPLOAD_COMMIT)

	err = upload.UploadCommit(args[1], metafile.Username)
	if err != nil {
		return fmt.Errorf("error uploading commit: %w", err)
	}
	fmt.Println("Commit uploaded successfuly!")
	return nil
}
