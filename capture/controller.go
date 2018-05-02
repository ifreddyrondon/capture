package capture

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"gopkg.in/src-d/go-kallax.v1"

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

// Controller handler the capture's routes
type Controller struct {
	Service Service
	render.Render
	CtxKey fmt.Stringer
}

// Router creates a REST router for the capture resource
func (c *Controller) Router() http.Handler {
	r := bastion.NewRouter()

	r.Get("/", c.list)
	r.Post("/", c.create)
	r.Route("/{id}", func(r chi.Router) {
		r.Use(c.captureCtx)
		r.Get("/", c.get)
		r.Put("/", c.update)
		r.Delete("/", c.delete)
	})
	return r
}

func (c *Controller) list(w http.ResponseWriter, r *http.Request) {
	count := 10
	start := 0
	captures, err := c.Service.List(start, count)
	if err != nil {
		_ = c.Render(w).InternalServerError(err)
		return
	}

	_ = c.Render(w).Send(captures)
}

func (c *Controller) create(w http.ResponseWriter, r *http.Request) {
	var captures Captures
	if err := json.NewDecoder(r.Body).Decode(&captures); err != nil {
		_ = c.Render(w).BadRequest(err)
		return
	}

	if len(captures) == 1 {
		if err := c.Service.Save(captures[0]); err != nil {
			_ = c.Render(w).InternalServerError(err)
			return
		}
		_ = c.Render(w).Created(captures[0])
		return
	}

	captures, err := c.Service.SaveBulk(captures...)
	if err != nil {
		_ = c.Render(w).InternalServerError(err)
		return
	}
	_ = c.Render(w).Created(captures)
}

func (c *Controller) captureCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captureID, err := kallax.NewULIDFromText(chi.URLParam(r, "id"))
		if err != nil {
			log.Println(err)
			_ = c.Render(w).BadRequest(ErrorBadRequest)
			return
		}
		var capt *Capture
		capt, err = c.Service.Get(captureID)
		if capt == nil {
			_ = c.Render(w).NotFound(err)
			return
		}
		if err != nil {
			_ = c.Render(w).InternalServerError(err)
			return
		}
		ctx := context.WithValue(r.Context(), c.CtxKey, capt)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (c *Controller) get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	capt, ok := ctx.Value(c.CtxKey).(*Capture)
	if !ok {
		err := errors.New(http.StatusText(http.StatusUnprocessableEntity))
		_ = c.Render(w).InternalServerError(err)
		return
	}
	_ = c.Render(w).Send(capt)
}

func (c *Controller) delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	capt, ok := ctx.Value(c.CtxKey).(*Capture)
	if !ok {
		err := errors.New(http.StatusText(http.StatusUnprocessableEntity))
		_ = c.Render(w).InternalServerError(err)
		return
	}
	if err := c.Service.Delete(capt); err != nil {
		_ = c.Render(w).InternalServerError(err)
		return
	}
	c.Render(w).NoContent()
}

func (c *Controller) update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	capt, ok := ctx.Value(c.CtxKey).(*Capture)
	if !ok {
		err := errors.New(http.StatusText(http.StatusUnprocessableEntity))
		_ = c.Render(w).InternalServerError(err)
		return
	}

	var updates Capture
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		_ = c.Render(w).BadRequest(err)
		return
	}

	if err := c.Service.Update(capt, updates); err != nil {
		_ = c.Render(w).InternalServerError(err)
		return
	}
	_ = c.Render(w).Send(capt)
}
