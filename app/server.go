package app

import (
	"fmt"

	"github.com/ifreddyrondon/gobastion"
)

type App struct {
	Bastion *gobastion.Bastion
}

func New(routes []Router) *App {
	server := &App{
		Bastion: gobastion.New(nil),
	}

	for _, v := range routes {
		server.Bastion.APIRouter.Mount(
			fmt.Sprintf("/%v/", v.Pattern()),
			v.Router(),
		)
	}

	return server
}
