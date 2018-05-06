package auth

import (
	"encoding/json"
	"net/http"

	"github.com/ifreddyrondon/bastion/render"

	"github.com/ifreddyrondon/bastion"
)

// Controller handler the auth routes
type Controller struct {
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

}
