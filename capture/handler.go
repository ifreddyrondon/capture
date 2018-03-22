package capture

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/render"
	bastionJSON "github.com/ifreddyrondon/bastion/render/json"
)

type Handler struct {
	Service Service
	render.Render
	CtxKey fmt.Stringer
}

// Router creates a REST router for the capture resource
func (h *Handler) Router() http.Handler {
	r := bastion.NewRouter()

	r.Get("/", h.list)
	r.Post("/", h.create)
	r.Route("/{id}", func(r chi.Router) {
		r.Use(h.captureCtx)
		r.Get("/", h.get)
		// r.Put("/", h.update)    // PUT /todos/{id} - update a single todo by :id
		// r.Delete("/", h.delete) // DELETE /todos/{id} - delete a single todo by :id
	})
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
}

func (h *Handler) captureCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captureID := chi.URLParam(r, "id")
		cap, err := h.Service.Get(captureID)
		if err != nil {
			h.Render(w).NotFound(err)
			return
		}
		ctx := context.WithValue(r.Context(), h.CtxKey, cap)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cap, ok := ctx.Value(h.CtxKey).(*Capture)
	if !ok {
		message := bastionJSON.HTTPError{
			Status: http.StatusUnprocessableEntity,
			Errors: http.StatusText(http.StatusUnprocessableEntity),
		}
		h.Render(w).Response(http.StatusUnprocessableEntity, message)
		return
	}
	h.Render(w).Send(cap)
}
