package main

import (
	"log"

	"github.com/ifreddyrondon/gocapture/app"
	"github.com/ifreddyrondon/gocapture/database"
)

func main() {
	ds, err := database.Open("localhost/captures")
	if err != nil {
		log.Panic(err)
	}

	err = app.New(ds.DB()).Serve()
	if err != nil {
		log.Printf("%v", err)
	}
}
