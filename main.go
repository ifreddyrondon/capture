package main

import (
	"log"

	"github.com/ifreddyrondon/gobastion"
	"github.com/ifreddyrondon/gocapture/branch"
	"github.com/ifreddyrondon/gocapture/capture"
	"github.com/ifreddyrondon/gocapture/database"
)

func main() {
	db, err := database.Open("localhost/captures")
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	app := gobastion.New("")

	app.APIRouter.Use(db.Ctx)
	app.APIRouter.Mount("/collection", new(branch.Handler).Routes())
	app.APIRouter.Mount("/captures", new(capture.Handler).Routes())
	err = app.Serve()
	if err != nil {
		log.Printf("%v", err)
	}
}
