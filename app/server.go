package app

import (
	"fmt"
	"log"

	"github.com/ifreddyrondon/gobastion"
	"github.com/ifreddyrondon/gocapture/branch"
	"github.com/ifreddyrondon/gocapture/capture"
	"github.com/ifreddyrondon/gocapture/database"
)

type App struct {
	Bastion *gobastion.Bastion
}

func Mount(routes []Router) *App {
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

func New() *App {
	ds, err := database.Open("localhost/captures")
	if err != nil {
		log.Panic(err)
	}

	reader := new(gobastion.JsonReader)
	captureService := capture.MgoService{DB: ds.DB()}
	captureHandler := capture.Handler{
		Service:   &captureService,
		Reader:    reader,
		Responder: gobastion.DefaultResponder,
	}

	branchHandler := branch.Handler{
		Reader:    reader,
		Responder: gobastion.DefaultResponder,
	}

	return Mount([]Router{&captureHandler, &branchHandler})
}
