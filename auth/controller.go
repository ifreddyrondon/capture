package auth

import (
	"errors"
	"net/http"

	"github.com/ifreddyrondon/gocapture/user"

	"github.com/go-chi/chi"
	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/render"
)

var errInvalidCredentials = errors.New("invalid email or password")

// Controller handler the auth routes
type Controller struct {
	render.Render
	*Strategies
}

// Router creates a REST router for the auth resource
func (c *Controller) Router() http.Handler {
	r := bastion.NewRouter()

	r.Route("/token-auth", func(r chi.Router) {
		r.Use(c.Strategies.LocalStrategy)
		r.Post("/", c.login)
	})
	return r
}

func (c *Controller) login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	u, ok := ctx.Value(c.Strategies.CtxKey).(*user.User)
	if !ok {
		err := errors.New(http.StatusText(http.StatusUnprocessableEntity))
		_ = c.Render(w).InternalServerError(err)
		return
	}
	_ = c.Render(w).Send(u)
}
