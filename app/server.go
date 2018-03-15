package app

import (
	"github.com/ifreddyrondon/bastion"
	"gopkg.in/mgo.v2"

	"github.com/ifreddyrondon/bastion/render/json"
	"github.com/ifreddyrondon/gocapture/branch"
	"github.com/ifreddyrondon/gocapture/capture"
)

func New(db *mgo.Database) *bastion.Bastion {
	app := bastion.New(bastion.Options{})

	capService := capture.MgoService{DB: db}
	capH := capture.Handler{
		Service: &capService,
		Render:  json.NewRender,
	}
	app.APIRouter.Mount(capH.Pattern(), capH.Router())

	branchH := branch.Handler{
		Render: json.NewRender,
	}
	app.APIRouter.Mount(branchH.Pattern(), branchH.Router())
	return app
}
