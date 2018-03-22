package branch

import (
	"encoding/json"
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
	r.Post("/", h.Create)
	return r
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	// the branch must be initialized to return an empty list when there are no captures.
	path := new(Branch)
	if err := json.NewDecoder(r.Body).Decode(path); err != nil {
		h.Render(w).BadRequest(err)
	}
	h.Render(w).Send(path)
	return
}
