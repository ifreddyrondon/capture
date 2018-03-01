package app

import (
	"net/http"

	"github.com/ifreddyrondon/bastion"
)

// Router is the interface implemented by the controllers.
// It allows the auto attach of the Router() (http.Handler)
// as a subrouter along a routing Pattern()
type Router interface {
	Pattern() string
	Router() http.Handler
}

type Handler interface {
	Router
	bastion.Reader
	bastion.Responder
}
