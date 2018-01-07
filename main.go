package main

import (
	"log"

	"github.com/ifreddyrondon/gocapture/database"
	"github.com/ifreddyrondon/gocapture/server"
)

func main() {
	db, err := database.Open("localhost/captures")
	if err != nil {
		log.Panic(err)
	}

	app := server.New(db)
	err = app.Bastion.Serve()
	if err != nil {
		log.Printf("%v", err)
	}
}
