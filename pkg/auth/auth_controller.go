package auth

import (
	"net/http"

	"github.com/ifreddyrondon/capture/pkg/auth/jwt"

	"github.com/ifreddyrondon/capture/pkg/user"

	"github.com/go-chi/chi"
	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/render"
)

type tokenJSON struct {
	Token string `json:"token,omitempty"`
}

// Routes returns a configured http.Handler with capture resources.
func Routes(authenticate func(http.Handler) http.Handler, service *jwt.Service) http.Handler {
	c := &controller{service: service, render: render.NewJSON()}

	r := bastion.NewRouter()
	r.Route("/token-auth", func(r chi.Router) {
		r.Use(authenticate)
		r.Post("/", c.login)
	})
	return r
}

type controller struct {
	render  render.APIRenderer
	service *jwt.Service
}

func (c *controller) login(w http.ResponseWriter, r *http.Request) {
	u, err := user.GetFromContext(r.Context())
	if err != nil {
		c.render.InternalServerError(w, err)
		return
	}

	token, err := c.service.GenerateToken(u.ID.String())
	if err != nil {
		c.render.InternalServerError(w, err)
	}

	c.render.Send(w, tokenJSON{Token: token})
}
