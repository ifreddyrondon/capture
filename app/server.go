package app

import (
	"fmt"

	"github.com/ifreddyrondon/gobastion"
	"github.com/ifreddyrondon/gocapture/database"
)

type App struct {
	Bastion *gobastion.Bastion
	*database.DB
}

func New(db *database.DB, routes []Router) *App {
	server := &App{
		Bastion: gobastion.New(nil),
		DB:      db,
	}
	server.Bastion.AppendFinalizers(db)
	server.Bastion.APIRouter.Use(db.Ctx)

	for _, v := range routes {
		server.Bastion.APIRouter.Mount(
			fmt.Sprintf("/%v/", v.Pattern()),
			v.Router(),
		)
	}

	return server
}
