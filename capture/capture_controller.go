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
	service Service
	render  render.Render
	ctxKey  fmt.Stringer
}

// NewController returns a new Controller
func NewController(service Service, render render.Render, ctxKey fmt.Stringer) *Controller {
	return &Controller{service: service, render: render, ctxKey: ctxKey}
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
	captures, err := c.service.List(start, count)
	if err != nil {
		_ = c.render(w).InternalServerError(err)
		return
	}

	_ = c.render(w).Send(captures)
}

func (c *Controller) create(w http.ResponseWriter, r *http.Request) {
	var captures Captures
	if err := json.NewDecoder(r.Body).Decode(&captures); err != nil {
		_ = c.render(w).BadRequest(err)
		return
	}

	if len(captures) == 1 {
		if err := c.service.Save(captures[0]); err != nil {
			_ = c.render(w).InternalServerError(err)
			return
		}
		_ = c.render(w).Created(captures[0])
		return
	}

	captures, err := c.service.SaveBulk(captures...)
	if err != nil {
		_ = c.render(w).InternalServerError(err)
		return
	}
	_ = c.render(w).Created(captures)
}

func (c *Controller) captureCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captureID, err := kallax.NewULIDFromText(chi.URLParam(r, "id"))
		if err != nil {
			log.Println(err)
			_ = c.render(w).BadRequest(ErrorBadRequest)
			return
		}
		var capt *Capture
		capt, err = c.service.Get(captureID)
		if capt == nil {
			_ = c.render(w).NotFound(err)
			return
		}
		if err != nil {
			_ = c.render(w).InternalServerError(err)
			return
		}
		ctx := context.WithValue(r.Context(), c.ctxKey, capt)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (c *Controller) get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	capt, ok := ctx.Value(c.ctxKey).(*Capture)
	if !ok {
		err := errors.New(http.StatusText(http.StatusUnprocessableEntity))
		_ = c.render(w).InternalServerError(err)
		return
	}
	_ = c.render(w).Send(capt)
}

func (c *Controller) delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	capt, ok := ctx.Value(c.ctxKey).(*Capture)
	if !ok {
		err := errors.New(http.StatusText(http.StatusUnprocessableEntity))
		_ = c.render(w).InternalServerError(err)
		return
	}
	if err := c.service.Delete(capt); err != nil {
		_ = c.render(w).InternalServerError(err)
		return
	}
	c.render(w).NoContent()
}

func (c *Controller) update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	capt, ok := ctx.Value(c.ctxKey).(*Capture)
	if !ok {
		err := errors.New(http.StatusText(http.StatusUnprocessableEntity))
		_ = c.render(w).InternalServerError(err)
		return
	}

	var updates Capture
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		_ = c.render(w).BadRequest(err)
		return
	}

	if err := c.service.Update(capt, updates); err != nil {
		_ = c.render(w).InternalServerError(err)
		return
	}
	_ = c.render(w).Send(capt)
}
