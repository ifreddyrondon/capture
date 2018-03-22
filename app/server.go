package app

import (
	"github.com/ifreddyrondon/bastion"
	"gopkg.in/mgo.v2"

	"github.com/ifreddyrondon/bastion/render/json"
	"github.com/ifreddyrondon/gocapture/branch"
	"github.com/ifreddyrondon/gocapture/capture"
)

// ContextKey represent a key in a context.ContextValue
type ContextKey string

func (c ContextKey) String() string {
	return string(c)
}

func New(db *mgo.Database) *bastion.Bastion {
	app := bastion.New(bastion.Options{})

	capService := capture.MgoService{Collection: db.C("captures")}
	capH := capture.Handler{
		Service: &capService,
		Render:  json.NewRender,
		CtxKey:  ContextKey("capture"),
	}
	app.APIRouter.Mount("/captures/", capH.Router())

	branchH := branch.Handler{
		Render: json.NewRender,
	}
	app.APIRouter.Mount("/branches/", branchH.Router())
	return app
}
