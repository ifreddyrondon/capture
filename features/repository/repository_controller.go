package repository

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/render"
	"github.com/ifreddyrondon/capture/features/auth/authorization"
	"github.com/ifreddyrondon/capture/features/user"
)

// Controller handler the repository routes
type Controller struct {
	service        Service
	render         render.APIRenderer
	authorization  *authorization.Authorization
	userService    user.GetterService
	ctxUserManager *user.ContextManager
}

// NewController returns a new Controller
func NewController(service Service, authMiddleware *authorization.Authorization, userService user.GetterService) *Controller {
	return &Controller{
		service:        service,
		render:         render.NewJSON(),
		authorization:  authMiddleware,
		userService:    userService,
		ctxUserManager: user.NewContextManager(),
	}
}

// Router creates a REST router for the user resource
func (c *Controller) Router() http.Handler {
	r := bastion.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Use(c.authorization.IsAuthorizedREQ)
		r.Use(user.LoggedUser(c.userService))
		r.Post("/", c.create)
	})

	return r
}

func (c *Controller) create(w http.ResponseWriter, r *http.Request) {
	var repo Repository
	if err := json.NewDecoder(r.Body).Decode(&repo); err != nil {
		c.render.BadRequest(w, err)
		return
	}

	userID, err := c.ctxUserManager.GetUserID(r.Context())
	if err != nil {
		c.render.InternalServerError(w, err)
		return
	}

	repo.UserID = userID

	if err := c.service.Save(&repo); err != nil {
		c.render.InternalServerError(w, err)
		return
	}

	c.render.Created(w, repo)
}
