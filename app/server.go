package app

import (
	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/render/json"
	"github.com/ifreddyrondon/gocapture/branch"
	"github.com/ifreddyrondon/gocapture/capture"
	"github.com/jinzhu/gorm"
)

// ContextKey represent a key in a context.ContextValue
type ContextKey string

func (c ContextKey) String() string {
	return string(c)
}

func New(db *gorm.DB) *bastion.Bastion {
	app := bastion.New(bastion.Options{})

	capService := capture.PGService{DB: db}
	capService.Drop()
	capService.Migrate()
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
