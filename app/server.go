package app

import (
	"fmt"

	"github.com/ifreddyrondon/gobastion"
)

type App struct {
	Bastion *gobastion.Bastion
	DataSource
}

func New(ds DataSource, routes []Router) *App {
	server := &App{
		Bastion:    gobastion.New(nil),
		DataSource: ds,
	}
	server.Bastion.AppendFinalizers(ds)
	server.Bastion.APIRouter.Use(ds.GetCtx())

	for _, v := range routes {
		server.Bastion.APIRouter.Mount(
			fmt.Sprintf("/%v/", v.Pattern()),
			v.Router(),
		)
	}

	return server
}
