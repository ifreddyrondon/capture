package main

import (
	"log"

	"github.com/ifreddyrondon/gobastion"
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

	reader := new(gobastion.JsonReader)

	captureHandler := new(capture.Handler)
	captureHandler.Reader = reader
	captureHandler.Responder = gobastion.DefaultResponder

	branchHandler := new(branch.Handler)
	branchHandler.Reader = reader
	branchHandler.Responder = gobastion.DefaultResponder

	routers := []app.Router{captureHandler, branchHandler}

	err = app.New(db, routers).Bastion.Serve()
	if err != nil {
		log.Printf("%v", err)
	}
}
