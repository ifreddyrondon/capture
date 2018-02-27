package app

import (
	"github.com/go-chi/chi"
)

type Router interface {
	Pattern() string
	Router() chi.Router
}
