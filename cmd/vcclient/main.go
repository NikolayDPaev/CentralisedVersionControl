package main

import (
	"os"

	"github.com/NikolayDPaev/CentralisedVersionControl/client/client"
)

func main() {
	client.ReadArgs(os.Args[1:])
}
