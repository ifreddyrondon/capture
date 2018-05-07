package auth

import (
	"encoding/json"
	"net/http"

	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/render"
	bastionJSON "github.com/ifreddyrondon/bastion/render/json"
)

// Controller handler the auth routes
type Controller struct {
	Service
	render.Render
}

// Router creates a REST router for the auth resource
func (c *Controller) Router() http.Handler {
	r := bastion.NewRouter()

	r.Post("/token-auth", c.login)
	return r
}

func (c *Controller) login(w http.ResponseWriter, r *http.Request) {
	var credentials BasicAuthCrendential
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		_ = c.Render(w).BadRequest(err)
		return
	}

	u, err := c.Service.Authenticate(&credentials)
	if err != nil {
		if err == errInvalidCredentials {
			httpErr := bastionJSON.HTTPError{
				Status:  http.StatusUnauthorized,
				Errors:  http.StatusText(http.StatusUnauthorized),
				Message: err.Error(),
			}
			_ = c.Render(w).Response(http.StatusUnauthorized, httpErr)
			return
		}

		_ = c.Render(w).InternalServerError(err)
		return
	}

	_ = c.Render(w).Created(u)
}
