package rest

import (
	"fmt"
	"net/http"
	"os"

	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/render"
	"github.com/ifreddyrondon/capture/pkg/signup"
	"github.com/pkg/errors"
)

type conflictErr interface {
	Conflict() bool
}

func isConflictErr(err error) bool {
	if e, ok := errors.Cause(err).(conflictErr); ok {
		return e.Conflict()
	}
	return false
}

// SignUp returns a configured http.Handler with sign-up resources.
func SignUp(service signup.Service) http.Handler {
	c := &signUpController{service: service, render: render.NewJSON()}

	r := bastion.NewRouter()
	r.Post("/", c.create)
	return r
}

type signUpController struct {
	service signup.Service
	render  render.APIRenderer
}

func (c *signUpController) create(w http.ResponseWriter, r *http.Request) {
	var payl signup.Payload
	err := signup.Validator.Decode(r, &payl)
	if err != nil {
		c.render.BadRequest(w, err)
		return
	}

	u, err := c.service.EnrollUser(payl)
	if err != nil {
		if isConflictErr(err) {
			httpErr := render.HTTPError{
				Status:  http.StatusConflict,
				Error:   http.StatusText(http.StatusConflict),
				Message: fmt.Sprintf("email '%v' already exists", *payl.Email),
			}
			c.render.Response(w, http.StatusConflict, httpErr)
			return
		}

		fmt.Fprintln(os.Stderr, err)
		c.render.InternalServerError(w, err)
		return
	}

	c.render.Created(w, u)
}
