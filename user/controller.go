package user

import (
	"encoding/json"
	"net/http"

	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/render"
)

// Controller handler the user's routes
type Controller struct {
	Service Service
	render.Render
}

// Router creates a REST router for the user resource
func (c *Controller) Router() http.Handler {
	r := bastion.NewRouter()

	r.Post("/", c.create)
	return r
}

func (c *Controller) create(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		_ = c.Render(w).BadRequest(err)
		return
	}

	if err := c.Service.Save(&user); err != nil {
		_ = c.Render(w).InternalServerError(err)
		return
	}

	_ = c.Render(w).Created(user)
}
