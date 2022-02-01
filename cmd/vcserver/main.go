package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/NikolayDPaev/CentralisedVersionControl/server/server"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Provide interface address and port!")
		return
	}

	server := server.NewServer(os.Args[1], os.Args[2])
	server.Start()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() && scanner.Text() != "exit" {
	}

	server.Stop()
}
