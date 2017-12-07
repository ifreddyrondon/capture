package main

import (
	"github.com/ifreddyrondon/gobastion"
	"github.com/ifreddyrondon/gocapture/branch"
)

func main() {
	app := gobastion.New("")
	app.APIRouter.Mount("/collection", new(branch.Handler).Routes())
	app.Serve()
}
