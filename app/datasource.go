package app

import (
	"net/http"
)

// DataSource is an interface to manage the diferent kinds of databases
// from the aplicacion environment
type DataSource interface {
	// GetCtx returns a function that set into the context request a value.
	// The value represent a DataSource.
	GetCtx() func(next http.Handler) http.Handler
	// Finalize is a func that will be executed into the graceful shutdown of bastion.
	Finalize() error
}
