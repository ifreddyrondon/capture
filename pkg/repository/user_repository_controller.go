package repository

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/middleware"
	"github.com/ifreddyrondon/bastion/middleware/listing/filtering"
	"github.com/ifreddyrondon/bastion/middleware/listing/sorting"
	"github.com/ifreddyrondon/bastion/render"
	"github.com/ifreddyrondon/capture/pkg/repository/decoder"
	"github.com/ifreddyrondon/capture/pkg/repository/encoder"
	"github.com/ifreddyrondon/capture/pkg/user"
)

var (
	// ErrorInvalidRepoID expected error when repository id param is invalid
	ErrorInvalidRepoID = errors.New("invalid repository id")
	// ErrorNotFound expected error when repository is missing
	ErrorNotFound = errors.New("not found repository")
	// ErrorNotFound expected error when repository is missing
	ErrorNotAuthorized = errors.New("not authorized to see this repository")
)

// UserRoutes returns a configured http.Handler with user repositories resources.
func UserRoutes(service Service, isAuth, loggedUser func(http.Handler) http.Handler) http.Handler {
	c := &userController{service: service, render: render.NewJSON()}

	updatedDESC := sorting.NewSort("updated_at_desc", "updated_at DESC", "Updated date descending")
	updatedASC := sorting.NewSort("updated_at_asc", "updated_at ASC", "Updated date ascendant")
	createdDESC := sorting.NewSort("created_at_desc", "created_at DESC", "Created date descending")
	createdASC := sorting.NewSort("created_at_asc", "created_at ASC", "Created date ascendant")

	publicVisibility := filtering.NewValue("public", "public repos")
	privateVisibility := filtering.NewValue("private", "private repos")
	visibilityFilter := filtering.NewText("visibility", "filters the repos by their visibility", publicVisibility, privateVisibility)

	listing := middleware.Listing(
		middleware.MaxAllowedLimit(50),
		middleware.Sort(updatedDESC, updatedASC, createdDESC, createdASC),
		middleware.Filter(visibilityFilter),
	)

	r := bastion.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Use(isAuth)
		r.Use(loggedUser)
		r.Post("/", c.create)
		r.Route("/", func(r chi.Router) {
			r.Use(listing)
			r.Get("/", c.list)
		})
	})

	return r
}

type userController struct {
	service Service
	render  render.APIRenderer
}

func (c *userController) create(w http.ResponseWriter, r *http.Request) {
	var postRepo decoder.PostRepository
	if err := decoder.Decode(r, &postRepo); err != nil {
		c.render.BadRequest(w, err)
		return
	}

	repo := postRepo.GetRepository()
	u, err := user.GetFromContext(r.Context())
	if err != nil {
		c.render.InternalServerError(w, err)
		return
	}

	if err := c.service.Save(u, &repo); err != nil {
		c.render.InternalServerError(w, err)
		return
	}

	c.render.Created(w, repo)
}

func (c *userController) list(w http.ResponseWriter, r *http.Request) {
	listing, err := middleware.GetListing(r.Context())
	if err != nil {
		c.render.InternalServerError(w, err)
		return
	}

	u, err := user.GetFromContext(r.Context())
	if err != nil {
		c.render.InternalServerError(w, err)
		return
	}

	repos, err := c.service.GetUserRepositories(u, listing)
	if err != nil {
		c.render.InternalServerError(w, err)
		return
	}

	res := encoder.ListRepositoryResponse{Listing: listing, Results: repos}
	c.render.Send(w, res)
}
