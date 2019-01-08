package handler

import (
	"fmt"
	"net/http"
	"os"

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
func Authenticating(service authenticating.Service) http.HandlerFunc {
	renderJSON := render.NewJSON()

	return func(w http.ResponseWriter, r *http.Request) {
		var credential authenticating.BasicCredential
		err := authenticating.Validator.Decode(r, &credential)
		if err != nil {
			renderJSON.BadRequest(w, err)
			return
		}

		u, err := service.AuthenticateUser(credential)
		if err != nil {
			if isInvalidCredential(err) {
				httpErr := render.HTTPError{
					Status:  http.StatusUnauthorized,
					Error:   http.StatusText(http.StatusUnauthorized),
					Message: "invalid email or password",
				}
				renderJSON.Response(w, http.StatusUnauthorized, httpErr)
				return
			}
			fmt.Fprintln(os.Stderr, err)
			renderJSON.InternalServerError(w, err)
			return
		}

		t, err := service.GetUserToken(u.ID)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			renderJSON.InternalServerError(w, err)
			return
		}

		renderJSON.Send(w, tokenJSON{Token: t})
	}
}
