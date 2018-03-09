package capture

import (
	"net/http"

	"encoding/json"

	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/render"
)

const Domain = "captures"

type Handler struct {
	Service Service
	render.Render
}

func (h *Handler) Pattern() string {
	return Domain
}

// Routes creates a REST router for the capture resource
func (h *Handler) Router() http.Handler {
	r := bastion.NewRouter()

	r.Get("/", h.list)
	r.Post("/", h.create)
	return r
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	count := 10
	start := 0
	captures, err := h.Service.List(start, count)
	if err != nil {
		h.Render(w).InternalServerError(err)
		return
	}

	h.Render(w).Send(captures)
	return
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	captureIn := new(Capture)
	if err := json.NewDecoder(r.Body).Decode(captureIn); err != nil {
		h.Render(w).BadRequest(err)
		return
	}

	captureOut, err := h.Service.Create(captureIn.Point, captureIn.Timestamp, captureIn.Payload)
	if err != nil {
		h.Render(w).InternalServerError(err)
		return
	}

	h.Render(w).Created(captureOut)
	return
}
