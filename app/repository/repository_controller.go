package repository

import (
	json "encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/render"
	"github.com/ifreddyrondon/capture/app/auth/authorization"
	"github.com/ifreddyrondon/capture/app/user"
)

// Controller handler the repository routes
type Controller struct {
	service        Service
	render         render.Render
	authorization  *authorization.Authorization
	userMiddleware *user.Middleware
}

// NewController returns a new Controller
func NewController(service Service, render render.Render, authMiddleware *authorization.Authorization, userMiddleware *user.Middleware) *Controller {
	return &Controller{
		service:        service,
		render:         render,
		authorization:  authMiddleware,
		userMiddleware: userMiddleware,
	}
}

// Router creates a REST router for the user resource
func (c *Controller) Router() http.Handler {
	r := bastion.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Use(c.authorization.IsAuthorizedREQ)
		r.Use(c.userMiddleware.LoggedUser)
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
