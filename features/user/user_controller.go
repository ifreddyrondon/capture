package user

import (
	"encoding/json"
	"net/http"

	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/render"
)

// Routes returns a configured http.Handler with user resources.
func Routes(service Service) http.Handler {
	c := &controller{service: service, render: render.NewJSON()}

	r := bastion.NewRouter()
	r.Post("/", c.create)
	return r
}

type controller struct {
	service Service
	render  render.APIRenderer
}

func (c *controller) create(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		c.render.BadRequest(w, err)
		return
	}

	if err := c.service.Save(&user); err != nil {
		if _, ok := err.(*emailDuplicateError); ok {
			httpErr := render.HTTPError{
				Status:  http.StatusConflict,
				Error:   http.StatusText(http.StatusConflict),
				Message: err.Error(),
			}
			c.render.Response(w, http.StatusConflict, httpErr)
			return
		}

		c.render.InternalServerError(w, err)
		return
	}

	c.render.Created(w, user)
}
