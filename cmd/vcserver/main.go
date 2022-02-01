package main

import (
	"bufio"
	"os"

	"github.com/NikolayDPaev/CentralisedVersionControl/server/server"
)

func main() {
	server := server.NewServer("4444")
	server.Start()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() && scanner.Text() != "exit" {
	}

	server.Stop()
}
