package auth

import (
	"errors"
	"net/http"

	"github.com/ifreddyrondon/capture/app/auth/authentication"
	"github.com/ifreddyrondon/capture/app/auth/jwt"

	"github.com/ifreddyrondon/capture/app/user"

	"github.com/go-chi/chi"
	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/render"
)

type tokenJSON struct {
	Token string `json:"token,omitempty"`
}

// Controller handler the auth routes
type Controller struct {
	middleware *authentication.Authentication
	render     render.Render
	service    *jwt.Service
	ctxManager *user.ContextManager
}

// NewController returns a new Controller
func NewController(middleware *authentication.Authentication, service *jwt.Service, render render.Render) *Controller {
	return &Controller{
		middleware: middleware,
		service:    service,
		render:     render,
		ctxManager: user.NewContextManager(),
	}
}

// Router creates a REST router for the auth resource
func (c *Controller) Router() http.Handler {
	r := bastion.NewRouter()

	r.Route("/token-auth", func(r chi.Router) {
		r.Use(c.middleware.Authenticate)
		r.Post("/", c.login)
	})
	return r
}

func (c *Controller) login(w http.ResponseWriter, r *http.Request) {
	u := c.ctxManager.Get(r.Context())
	if u == nil {
		err := errors.New(http.StatusText(http.StatusUnprocessableEntity))
		_ = c.render(w).InternalServerError(err)
		return
	}

	token, err := c.service.GenerateToken(u.ID.String())
	if err != nil {
		_ = c.render(w).InternalServerError(err)
	}

	_ = c.render(w).Send(tokenJSON{Token: token})
}
