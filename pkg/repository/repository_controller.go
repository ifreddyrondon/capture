package repository

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/middleware"
	"github.com/ifreddyrondon/bastion/middleware/listing/sorting"
	"github.com/ifreddyrondon/bastion/render"
	"github.com/ifreddyrondon/capture/pkg/repository/encoder"
	"github.com/ifreddyrondon/capture/pkg/user"
	"gopkg.in/src-d/go-kallax.v1"
)

// Routes returns a configured http.Handler with repositories resources.
func Routes(service Service, isAuth, loggedUser func(http.Handler) http.Handler) http.Handler {
	c := &controller{service: service, render: render.NewJSON()}

	updatedDESC := sorting.NewSort("updated_at_desc", "updated_at DESC", "Updated date descending")
	updatedASC := sorting.NewSort("updated_at_asc", "updated_at ASC", "Updated date ascendant")
	createdDESC := sorting.NewSort("created_at_desc", "created_at DESC", "Created date descending")
	createdASC := sorting.NewSort("created_at_asc", "created_at ASC", "Created date ascendant")

	listing := middleware.Listing(
		middleware.MaxAllowedLimit(50),
		middleware.Sort(updatedDESC, updatedASC, createdDESC, createdASC),
	)

	r := bastion.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Use(listing)
		r.Get("/", c.list)
	})
	r.Route("/{id}", func(r chi.Router) {
		r.Use(isAuth)
		r.Use(loggedUser)
		r.Use(c.repoCtx)
		r.Get("/", c.get)
	})

	return r
}

type controller struct {
	service Service
	render  render.APIRenderer
}

func (c *controller) list(w http.ResponseWriter, r *http.Request) {
	listing, err := middleware.GetListing(r.Context())
	if err != nil {
		c.render.InternalServerError(w, err)
		return
	}

	repos, err := c.service.GetPublicRepositories(listing)
	if err != nil {
		c.render.InternalServerError(w, err)
		return
	}

	res := encoder.ListRepositoryResponse{Listing: listing, Results: repos}
	c.render.Send(w, res)
}

func (c *controller) repoCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		repoID, err := kallax.NewULIDFromText(chi.URLParam(r, "id"))
		if err != nil {
			c.render.BadRequest(w, ErrorInvalidRepoID)
			return
		}

		u, err := user.GetFromContext(r.Context())
		if err != nil {
			c.render.InternalServerError(w, err)
			return
		}

		repo, err := c.service.GetRepo(repoID, u)
		if err != nil {
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
