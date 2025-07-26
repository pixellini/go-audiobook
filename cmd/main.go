package main

import (
	"log"
	"os"

	"github.com/pixellini/go-audiobook/internal/cli"
)

func main() {
	err := cli.Run()

	if err != nil {
		log.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
