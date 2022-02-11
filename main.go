package main

import (
	"fmt"

	"github.com/NikolayDPaev/CentralisedVersionControl/client/fileio"
)

func main() {
	lf := &fileio.Localfiles{}
	// var bytebuff []byte
	// buf := bytes.NewBuffer(bytebuff)

	// comm := netio.NewCommunicator(10, buf, buf)

	// err := lf.SendBlob("test.txt", comm)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// err = lf.ReceiveBlob("test2.txt", comm)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	slice, _ := lf.GetPathsOfAllFiles()
	for i := 0; i < len(slice)-1; i++ {
		if slice[i] < slice[i+1] {
			fmt.Println("Not ordered")
		}
	}
}
