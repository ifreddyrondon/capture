package main

import (
	"log"

	"github.com/ifreddyrondon/gocapture/app"
)

func main() {
	err := app.New().Bastion.Serve()
	if err != nil {
		log.Printf("%v", err)
	}
}
