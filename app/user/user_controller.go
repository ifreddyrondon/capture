package user

import (
	"encoding/json"
	"net/http"

	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/render"
)

// Controller handler the user's routes
type Controller struct {
	service Service
	render  render.APIRenderer
}

// NewController returns a new Controller
func NewController(service Service) *Controller {
	return &Controller{service: service, render: render.NewJSON()}
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
