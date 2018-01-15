package branch

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/ifreddyrondon/gobastion"
)

const Domain = "branches"

type Handler struct {
	gobastion.Reader
	gobastion.Responder
}

func (h *Handler) Pattern() string {
	return Domain
}

// Routes creates a REST router for the branch resource
func (h *Handler) Router() chi.Router {
	r := gobastion.NewRouter()
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
