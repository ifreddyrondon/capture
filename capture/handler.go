package capture

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/ifreddyrondon/gobastion"
)

const Domain = "captures"

type Handler struct {
	Service Service
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
	count := 10
	start := 0
	captures, err := h.Service.List(start, count)
	if err != nil {
		h.InternalServerError(w, err)
		return
	}

	h.Send(w, captures)
	return
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	captureIn := new(Capture)
	if err := h.Read(r.Body, captureIn); err != nil {
		h.BadRequest(w, err)
		return
	}

	captureOut, err := h.Service.Create(captureIn.Point, captureIn.Timestamp, captureIn.Payload)
	if err != nil {
		h.InternalServerError(w, err)
		return
	}

	h.Created(w, captureOut)
	return
}
