package app

import (
	"net/http"

	"github.com/go-chi/chi"
)

type Router interface {
	Pattern() string
	Router() chi.Router
}

// DataSource is an interface to manage the different kinds of databases
// from the application environment
type DataSource interface {
	// Ctx returns a function that set into the context request a value.
	// The value represent a DataSource.
	Ctx() func(next http.Handler) http.Handler
	// Finalize will be executed into the graceful shutdown of bastion.
	Finalize() error
}
