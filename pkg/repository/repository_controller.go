package repository

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/render"
	auth "github.com/ifreddyrondon/capture/pkg/http/rest/middleware"
)

var (
	// ErrorNotFound expected error when repository is missing
	ErrorNotFound = errors.New("not found repository")
	// ErrorNotFound expected error when repository is missing
	ErrorNotAuthorized = errors.New("not authorized to see this repository")
)

// Routes returns a configured http.Handler with repositories resources.
func Routes(service Service, authorizeReq func(http.Handler) http.Handler) http.Handler {
	c := &controller{service: service, render: render.NewJSON()}

	r := bastion.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Use(authorizeReq)
		r.Use(c.repoCtx)
		r.Get("/", c.get)
	})

	return r
}

type controller struct {
	service Service
	render  render.APIRenderer
}

func (c *controller) repoCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		repoID := chi.URLParam(r, "id")
		u, err := auth.GetUser(r.Context())
		if err != nil {
			c.render.InternalServerError(w, err)
			return
		}

		repo, err := c.service.GetRepo(repoID, u)
		if err != nil {
			// FIXME: handler bad repo id should be BAD REQUEST
			if err == ErrorNotFound {
				c.render.NotFound(w, err)
				return
			}
			if err == ErrorNotAuthorized {
				s := http.StatusForbidden
				message := render.NewHTTPError(err.Error(), http.StatusText(s), s)
				c.render.Response(w, s, message)
				return
			}
			c.render.InternalServerError(w, err)
			return
		}
		ctx := withRepo(r.Context(), repo)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (c *controller) get(w http.ResponseWriter, r *http.Request) {
	repo, err := GetFromContext(r.Context())
	if err != nil {
		c.render.InternalServerError(w, err)
		return
	}
	c.render.Send(w, repo)
}
