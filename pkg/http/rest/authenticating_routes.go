package rest

import (
	"fmt"
	"net/http"
	"os"

	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/render"
	"github.com/ifreddyrondon/capture/pkg/authenticating"
	"github.com/pkg/errors"
)

type authenticatingErr interface {
	// InvalidCredentials returns true when the error is by invalid credentials
	InvalidCredentials() bool
}

func isInvalidCredential(err error) bool {
	if e, ok := errors.Cause(err).(authenticatingErr); ok {
		return e.InvalidCredentials()
	}
	return false
}

type tokenJSON struct {
	Token string `json:"token,omitempty"`
}

// AuthenticatingRoutes returns a configured http.Handler with capture resources.
func Authenticating(service authenticating.Service) http.Handler {
	c := &authenticatingController{service: service, render: render.NewJSON()}

	r := bastion.NewRouter()
	r.Post("/token-auth", c.login)
	return r
}

type authenticatingController struct {
	render  render.APIRenderer
	service authenticating.Service
}

func (c *authenticatingController) login(w http.ResponseWriter, r *http.Request) {
	var credential authenticating.BasicCredential
	err := authenticating.Validator.Decode(r, &credential)
	if err != nil {
		c.render.BadRequest(w, err)
		return
	}

	u, err := c.service.AuthenticateUser(credential)
	if err != nil {
		if isInvalidCredential(err) {
			httpErr := render.HTTPError{
				Status:  http.StatusUnauthorized,
				Error:   http.StatusText(http.StatusUnauthorized),
				Message: "invalid email or password",
			}
			c.render.Response(w, http.StatusUnauthorized, httpErr)
			return
		}
		fmt.Fprintln(os.Stderr, err)
		c.render.InternalServerError(w, err)
		return
	}

	t, err := c.service.GetUserToken(u.ID)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		c.render.InternalServerError(w, err)
		return
	}

	c.render.Send(w, tokenJSON{Token: t})
}
