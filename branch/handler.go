package branch

import (
	"net/http"

	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/render"
)

type Handler struct {
	render.Render
}

// Routes creates a REST router for the branch resource
func (h *Handler) Router() http.Handler {
	r := bastion.NewRouter()
	r.Get("/", h.list)
	return r
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	b := Branch{}
	h.Render(w).Send(b)
}
