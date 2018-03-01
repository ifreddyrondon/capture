package branch

import (
	"net/http"

	"github.com/ifreddyrondon/bastion"
)

const Domain = "branches"

type Handler struct {
	bastion.Reader
	bastion.Responder
}

func (h *Handler) Pattern() string {
	return Domain
}

// Routes creates a REST router for the branch resource
func (h *Handler) Router() http.Handler {
	r := bastion.NewRouter()
	r.Post("/", h.Create)
	return r
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	path := new(Branch)
	if err := h.Read(r.Body, path); err != nil {
		h.BadRequest(w, err)
	}
	h.Send(w, path)
	return
}
