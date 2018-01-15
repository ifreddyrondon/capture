package capture

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/ifreddyrondon/gobastion"
	"gopkg.in/mgo.v2"
)

const Domain = "captures"

type Handler struct {
	gobastion.Reader
	gobastion.Responder
}

func (h *Handler) Pattern() string {
	return Domain
}

// Routes creates a REST router for the capture resource
func (h *Handler) Router() chi.Router {
	r := gobastion.NewRouter()
	r.Post("/", h.create)
	return r
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	capture := new(Capture)
	if err := h.Read(r.Body, capture); err != nil {
		h.BadRequest(w, err)
		return
	}

	ctx := r.Context()
	err := capture.create(ctx.Value("DataSource").(*mgo.Database))
	if err != nil {
		h.InternalServerError(w, err)
		return
	}

	h.Created(w, capture)
	return
}
