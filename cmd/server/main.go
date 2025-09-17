package main

import (
	"log"

	"github.com/pixellini/go-audiobook/internal/server"
)

func main() {
	s := server.New()
	log.Fatal(s.Start(":8003"))
}
