package server

import (
	"github.com/ifreddyrondon/gobastion"
	"github.com/ifreddyrondon/gocapture/branch"
	"github.com/ifreddyrondon/gocapture/capture"
	"github.com/ifreddyrondon/gocapture/database"
)

type Server struct {
	Bastion *gobastion.Bastion
	*database.DB
}

func New(db *database.DB) *Server {
	server := &Server{
		Bastion: gobastion.New(nil),
		DB:      db,
	}
	server.Bastion.AppendFinalizers(db)
	server.Bastion.APIRouter.Use(db.Ctx)
	server.Bastion.APIRouter.Mount("/collection", new(branch.Handler).Routes())
	server.Bastion.APIRouter.Mount("/captures", new(capture.Handler).Routes())
	return server
}
