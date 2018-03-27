package capture

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/render"
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
		r.Put("/", h.update)
		r.Delete("/", h.delete)
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
	var captureIN Capture
	if err := json.NewDecoder(r.Body).Decode(&captureIN); err != nil {
		h.Render(w).BadRequest(err)
		return
	}

	captureOUT, err := h.Service.Create(captureIN.Point, captureIN.Timestamp, captureIN.Payload)
	if err != nil {
		h.Render(w).InternalServerError(err)
		return
	}

	h.Render(w).Created(captureOUT)
}

func (h *Handler) captureCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captureID := chi.URLParam(r, "id")
		cap, err := h.Service.Get(captureID)
		if cap == nil {
			h.Render(w).NotFound(err)
			return
		}
		if err != nil {
			h.Render(w).InternalServerError(err)
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
		err := errors.New(http.StatusText(http.StatusUnprocessableEntity))
		h.Render(w).InternalServerError(err)
		return
	}
	h.Render(w).Send(cap)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cap, ok := ctx.Value(h.CtxKey).(*Capture)
	if !ok {
		err := errors.New(http.StatusText(http.StatusUnprocessableEntity))
		h.Render(w).InternalServerError(err)
		return
	}
	if err := h.Service.Delete(cap.ID.Hex()); err != nil {
		h.Render(w).InternalServerError(err)
		return
	}
	h.Render(w).NoContent()
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cap, ok := ctx.Value(h.CtxKey).(*Capture)
	if !ok {
		err := errors.New(http.StatusText(http.StatusUnprocessableEntity))
		h.Render(w).InternalServerError(err)
		return
	}

	var captureDST Capture
	if err := json.NewDecoder(r.Body).Decode(&captureDST); err != nil {
		h.Render(w).BadRequest(err)
		return
	}

	captureDST.ID = cap.ID
	captureDST.Visible = cap.Visible
	captureDST.CreatedDate = cap.CreatedDate
	captureDST.LastModified = time.Now()

	err := h.Service.Update(&captureDST)
	if err != nil {
		h.Render(w).InternalServerError(err)
		return
	}
	h.Render(w).Send(captureDST)
}
