package capture

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/render"
)

var (
	// ErrorNotFound expected error when capture is missing
	ErrorNotFound = errors.New("not found capture")
	// ErrorBadRequest expected error when capture id is invalid
	ErrorBadRequest = errors.New("invalid capture id")
)

// Controller handler all the router for capture
type Controller struct {
	Service Service
	render.Render
	CtxKey fmt.Stringer
}

// Router creates a REST router for the capture resource
func (h *Controller) Router() http.Handler {
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

func (h *Controller) list(w http.ResponseWriter, r *http.Request) {
	count := 10
	start := 0
	captures, err := h.Service.List(start, count)
	if err != nil {
		_ = h.Render(w).InternalServerError(err)
		return
	}

	_ = h.Render(w).Send(captures)
}

func (h *Controller) create(w http.ResponseWriter, r *http.Request) {
	var captures Captures
	if err := json.NewDecoder(r.Body).Decode(&captures); err != nil {
		_ = h.Render(w).BadRequest(err)
		return
	}

	if len(captures) == 1 {
		capt, err := h.Service.Save(captures[0])
		if err != nil {
			_ = h.Render(w).InternalServerError(err)
			return
		}
		_ = h.Render(w).Created(capt)
		return
	}

	captures, err := h.Service.SaveBulk(captures...)
	if err != nil {
		_ = h.Render(w).InternalServerError(err)
		return
	}
	_ = h.Render(w).Created(captures)
}

func (h *Controller) captureCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captureID, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			log.Println(err)
			_ = h.Render(w).BadRequest(ErrorBadRequest)
			return
		}
		var capt *Capture
		capt, err = h.Service.Get(captureID)
		if capt == nil {
			_ = h.Render(w).NotFound(err)
			return
		}
		if err != nil {
			_ = h.Render(w).InternalServerError(err)
			return
		}
		ctx := context.WithValue(r.Context(), h.CtxKey, capt)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Controller) get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	capt, ok := ctx.Value(h.CtxKey).(*Capture)
	if !ok {
		err := errors.New(http.StatusText(http.StatusUnprocessableEntity))
		_ = h.Render(w).InternalServerError(err)
		return
	}
	_ = h.Render(w).Send(capt)
}

func (h *Controller) delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	capt, ok := ctx.Value(h.CtxKey).(*Capture)
	if !ok {
		err := errors.New(http.StatusText(http.StatusUnprocessableEntity))
		_ = h.Render(w).InternalServerError(err)
		return
	}
	if err := h.Service.Delete(capt); err != nil {
		_ = h.Render(w).InternalServerError(err)
		return
	}
	h.Render(w).NoContent()
}

func (h *Controller) update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	capt, ok := ctx.Value(h.CtxKey).(*Capture)
	if !ok {
		err := errors.New(http.StatusText(http.StatusUnprocessableEntity))
		_ = h.Render(w).InternalServerError(err)
		return
	}

	var updates Capture
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		_ = h.Render(w).BadRequest(err)
		return
	}

	if err := h.Service.Update(capt, updates); err != nil {
		_ = h.Render(w).InternalServerError(err)
		return
	}
	_ = h.Render(w).Send(capt)
}
