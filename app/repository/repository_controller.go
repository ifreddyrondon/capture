package repository

import (
	"encoding/json"
	"net/http"

	"github.com/ifreddyrondon/capture/app/auth/authorization"

	"github.com/go-chi/chi"

	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/render"
)

// Controller handler the repository routes
type Controller struct {
	service       Service
	render        render.Render
	authorization authorization.Authorization
}

// NewController returns a new Controller
func NewController(service Service, render render.Render, authMiddleware authorization.Authorization) *Controller {
	return &Controller{
		service:       service,
		render:        render,
		authorization: authMiddleware,
	}
}

// Router creates a REST router for the user resource
func (c *Controller) Router() http.Handler {
	r := bastion.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Use(c.authorization.IsAuthorized)
		r.Post("/", c.create)
	})

	return r
}

func (c *Controller) create(w http.ResponseWriter, r *http.Request) {
	var repo Repository
	if err := json.NewDecoder(r.Body).Decode(&repo); err != nil {
		_ = c.render(w).BadRequest(err)
		return
	}

	if err := c.service.Save(&repo); err != nil {
		_ = c.render(w).InternalServerError(err)
		return
	}

	_ = c.render(w).Created(repo)
}
