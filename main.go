package main

import (
	"github.com/ifreddyrondon/gobastion"
	"github.com/ifreddyrondon/gocapture/capture"
)

func main() {
	app := gobastion.New("")
	app.APIRouter.Mount("/collection", new(capture.Handler).Routes())
	app.Serve()
}
