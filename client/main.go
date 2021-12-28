package main

import (
	"fmt"
	"log"

	"github.com/NikolayDPaev/CentralisedVersionControl/client/fileIO"
)

func main() {
	// file, err := fileIO.CompressToTempFile("../../test.mp3")
	// if err != nil {
	// 	log.Println(err)
	// }
	// err = fileIO.DecompressFile("test.mp3", file)
	// if err != nil {
	// 	log.Println(err)
	// }

	slice, err := fileIO.GetPathsOfAllFiles()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(slice)
}
