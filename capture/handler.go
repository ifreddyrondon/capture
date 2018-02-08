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

	r.Get("/", h.list)
	r.Post("/", h.create)
	return r
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	capture := new(Capture)

	count := 10
	start := 0
	ctx := r.Context()
	captures, err := capture.list(ctx.Value("DataSource").(*mgo.Database), start, count)
	if err != nil {
		h.InternalServerError(w, err)
		return
	}

	h.Send(w, captures)
	return
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
