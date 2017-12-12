package main

import (
	"github.com/ifreddyrondon/gobastion"
	"github.com/ifreddyrondon/gocapture/branch"
	"github.com/ifreddyrondon/gocapture/capture"
	"github.com/ifreddyrondon/gocapture/database"
)

func main() {
	app := gobastion.New("")
	database.CreateConnection("localhost/captures")
	defer database.DB.Close()
	app.APIRouter.Use(database.MongoCtx)
	app.APIRouter.Mount("/collection", new(branch.Handler).Routes())
	app.APIRouter.Mount("/captures", new(capture.Handler).Routes())
	app.Serve()
}
