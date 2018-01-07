package main

import (
	"log"

	"github.com/ifreddyrondon/gocapture/app"
	"github.com/ifreddyrondon/gocapture/branch"
	"github.com/ifreddyrondon/gocapture/capture"
	"github.com/ifreddyrondon/gocapture/database"
)

func main() {
	db, err := database.Open("localhost/captures")
	if err != nil {
		log.Panic(err)
	}

	routers := []app.Router{
		new(capture.Handler),
		new(branch.Handler),
	}

	err = app.New(db, routers).Bastion.Serve()
	if err != nil {
		log.Printf("%v", err)
	}
}
