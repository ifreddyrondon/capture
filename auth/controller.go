package auth

import (
	"net/http"

	"github.com/ifreddyrondon/bastion"
)

// Controller handler the auth routes
type Controller struct{}

// Router creates a REST router for the auth resource
func (c *Controller) Router() http.Handler {
	r := bastion.NewRouter()

	//r.Post("/token-auth", c.login)
	return r
}

// func (c *Controller) login(w http.ResponseWriter, r *http.Request) {

// }
