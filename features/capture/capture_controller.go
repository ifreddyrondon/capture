package capture

import (
	"errors"
	"log"
	"net/http"

	"github.com/ifreddyrondon/capture/features"
	"github.com/ifreddyrondon/capture/features/capture/decoder"
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

// Routes returns a configured http.Handler with capture resources.
func Routes(service Service) http.Handler {
	c := &controller{service: service, render: render.NewJSON()}

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

// Controller handler the capture's routes
type controller struct {
	service Service
	render  render.APIRenderer
}

func (c *controller) list(w http.ResponseWriter, r *http.Request) {
	count := 10
	start := 0
	captures, err := c.service.List(start, count)
	if err != nil {
		c.render.InternalServerError(w, err)
		return
	}

	c.render.Send(w, captures)
}

func (c *controller) create(w http.ResponseWriter, r *http.Request) {
	var postCapture decoder.POSTCapture
	if err := decoder.Decode(r, &postCapture); err != nil {
		c.render.BadRequest(w, err)
		return
	}

	capt := postCapture.GetCapture()
	if err := c.service.Save(&capt); err != nil {
		c.render.InternalServerError(w, err)
		return
	}
	c.render.Created(w, capt)
	return
}

func (c *controller) captureCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captureID, err := kallax.NewULIDFromText(chi.URLParam(r, "id"))
		if err != nil {
			log.Println(w, err)
			c.render.BadRequest(w, ErrorBadRequest)
			return
		}
		var capt *features.Capture
		capt, err = c.service.Get(captureID)
		if capt == nil {
			c.render.NotFound(w, err)
			return
		}
		if err != nil {
			c.render.InternalServerError(w, err)
			return
		}
		ctx := WithCapture(r.Context(), capt)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (c *controller) get(w http.ResponseWriter, r *http.Request) {
	capt, err := GetFromContext(r.Context())
	if err != nil {
		c.render.InternalServerError(w, err)
		return
	}
	c.render.Send(w, capt)
}

func (c *controller) delete(w http.ResponseWriter, r *http.Request) {
	capt, err := GetFromContext(r.Context())
	if err != nil {
		c.render.InternalServerError(w, err)
		return
	}
	if err := c.service.Delete(capt); err != nil {
		c.render.InternalServerError(w, err)
		return
	}
	c.render.NoContent(w)
}

func (c *controller) update(w http.ResponseWriter, r *http.Request) {
	captFromCtx, err := GetFromContext(r.Context())
	if err != nil {
		c.render.InternalServerError(w, err)
		return
	}

	var putCapture decoder.PUTCapture
	if err := decoder.Decode(r, &putCapture); err != nil {
		c.render.BadRequest(w, err)
		return
	}

	captFromPayload := putCapture.GetCapture()
	if err := c.service.Update(captFromCtx, captFromPayload); err != nil {
		c.render.InternalServerError(w, err)
		return
	}
	c.render.Send(w, captFromCtx)
}
