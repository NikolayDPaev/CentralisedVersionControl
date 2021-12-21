package files

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"sort"

	"github.com/NikolayDPaev/CentralisedVersionControl/server/commit"
)

func extractMetadata(fileInfo fs.FileInfo) (*commit.Metadata, error) {
	file, err := os.Open("commits/" + fileInfo.Name())
	if err != nil {
		return nil, err
	}
	defer file.Close()

	result, err := commit.NewMetadataFromFile(file, fileInfo.Name())
	if err != nil {
		return nil, err
	}
	return result, nil
}

func GetCommitList() []*commit.Metadata {
	f, err := os.Open("./commits")
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	defer f.Close()

	files, err := f.Readdir(0)
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].ModTime().Before(files[j].ModTime())
	})

	result := make([]*commit.Metadata, len(files))
	for i, v := range files {
		result[i], err = extractMetadata(v)
		if err != nil {
			continue
		}
	}
	return result
}

func GetCommitReader(commitId string) (io.Reader, error) {
	reader, err := os.Open("commits/" + commitId)
	if err != nil {
		return nil, fmt.Errorf("cannot open commit: %w", err)
	}
	return reader, nil
}

func NewCommitWriter(commitId string) (io.Writer, error) {
	writer, err := os.OpenFile("commits/"+commitId, os.O_WRONLY|os.O_CREATE|os.O_EXCL|os.O_APPEND, 0660)
	if err != nil {
		return nil, fmt.Errorf("cannot create commit file: %w", err)
	}
	return writer, nil
}
