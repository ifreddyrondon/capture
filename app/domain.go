package app

import (
	"github.com/go-chi/chi"
	"github.com/ifreddyrondon/gobastion"
)

// Router is the interface implemented by the controllers.
// It allows the auto attach of the Router() (http.Handler)
// as a subrouter along a routing Pattern()
type Router interface {
	Pattern() string
	Router() chi.Router
}

type Handler interface {
	Router
	gobastion.Reader
	gobastion.Responder
}
