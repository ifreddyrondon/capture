package app

import (
	"net/http"
)

// Router is the interface implemented by the controllers.
// It allows the auto attach of the Router() (http.Handler)
// as a subrouter along a routing path
type Router interface {
	Router() http.Handler
}
