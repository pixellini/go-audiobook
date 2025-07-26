package main

import (
	"fmt"

	"github.com/pixellini/go-audiobook/internal/app"
)

func main() {
	// err := cli.Execute()

	// if err != nil {
	// 	log.Printf("Error: %v\n", err)
	// 	os.Exit(1)
	// }

	app, err := app.New()
	if err != nil {
		fmt.Println("error happened on create", err)
	}

	err = app.Run()
	if err != nil {
		fmt.Println("error happened on run", err)
	}
}
