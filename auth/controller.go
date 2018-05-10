package auth

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ifreddyrondon/gocapture/jwt"

	"github.com/ifreddyrondon/gocapture/auth/strategy/basic"
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
	basic.Strategy
	UserKey fmt.Stringer
	render.Render
	JWT *jwt.Service
}

// Router creates a REST router for the auth resource
func (c *Controller) Router() http.Handler {
	r := bastion.NewRouter()

	r.Route("/token-auth", func(r chi.Router) {
		r.Use(c.Strategy.Authenticate)
		r.Post("/", c.login)
	})
	return r
}

func (c *Controller) login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	u, ok := ctx.Value(c.UserKey).(*user.User)
	if !ok {
		err := errors.New(http.StatusText(http.StatusUnprocessableEntity))
		_ = c.Render(w).InternalServerError(err)
		return
	}

	token, err := c.JWT.GenerateToken(u.ID.String())
	if err != nil {
		_ = c.Render(w).InternalServerError(err)
	}

	_ = c.Render(w).Send(tokenJSON{Token: token})
}
