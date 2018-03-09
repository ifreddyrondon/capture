package app

import (
	"fmt"
	"log"

	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/render/json"
	"github.com/ifreddyrondon/gocapture/branch"
	"github.com/ifreddyrondon/gocapture/capture"
	"github.com/ifreddyrondon/gocapture/database"
)

type App struct {
	Bastion *bastion.Bastion
}

func Mount(routes []Router) *App {
	server := &App{
		Bastion: bastion.New(nil),
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

	captureService := capture.MgoService{DB: ds.DB()}
	captureHandler := capture.Handler{
		Service: &captureService,
		Render:  json.NewRender,
	}

	branchHandler := branch.Handler{
		Render: json.NewRender,
	}

	return Mount([]Router{&captureHandler, &branchHandler})
}
