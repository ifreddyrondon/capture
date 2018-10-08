package user

import (
	"net/http"

	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/render"
	"github.com/ifreddyrondon/capture/features"
	"github.com/ifreddyrondon/capture/features/user/decoder"
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
	var postUser decoder.PostUser
	if err := decoder.Decode(r, &postUser); err != nil {
		c.render.BadRequest(w, err)
		return
	}

	var u features.User
	if err := decoder.User(postUser, &u); err != nil {
		c.render.InternalServerError(w, err)
		return
	}

	if err := c.service.Save(&u); err != nil {
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

	c.render.Created(w, u)
}
