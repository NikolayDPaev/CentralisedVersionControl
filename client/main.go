package main

import (
	"os"

	"github.com/NikolayDPaev/CentralisedVersionControl/client/clientUI"
)

func main() {
	clientUI.ReadArgs(os.Args[1:])
}
