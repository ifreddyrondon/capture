package auth

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ifreddyrondon/gocapture/jwt"
	"github.com/ifreddyrondon/gocapture/user"

	"github.com/go-chi/chi"
	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/render"
)

type tokenJSON struct {
	Token string `json:"token,omitempty"`
}

// Controller handler the auth routes
type Controller struct {
	strategy Strategy
	userKey  fmt.Stringer
	render   render.Render
	service  *jwt.Service
}

// NewController returns a new Controller
func NewController(strategy Strategy, service *jwt.Service, render render.Render, userKey fmt.Stringer) *Controller {
	return &Controller{
		strategy: strategy,
		service:  service,
		render:   render,
		userKey:  userKey,
	}
}

// Router creates a REST router for the auth resource
func (c *Controller) Router() http.Handler {
	r := bastion.NewRouter()

	r.Route("/token-auth", func(r chi.Router) {
		r.Use(c.strategy.Authenticate)
		r.Post("/", c.login)
	})
	return r
}

func (c *Controller) login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	u, ok := ctx.Value(c.userKey).(*user.User)
	if !ok {
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
