package handler

import (
	"fmt"
	"net/http"
	"os"

	"github.com/ifreddyrondon/bastion/binder"
	"github.com/ifreddyrondon/bastion/render"
	"github.com/pkg/errors"

	"github.com/ifreddyrondon/capture/pkg/authenticating"
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
	return func(w http.ResponseWriter, r *http.Request) {
		var credential authenticating.BasicCredential
		if err := binder.JSON.FromReq(r, &credential); err != nil {
			render.JSON.BadRequest(w, err)
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
				render.JSON.Response(w, http.StatusUnauthorized, httpErr)
				return
			}
			fmt.Fprintln(os.Stderr, err)
			render.JSON.InternalServerError(w, err)
			return
		}

		t, err := service.GetUserToken(u.ID)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			render.JSON.InternalServerError(w, err)
			return
		}

		render.JSON.Send(w, tokenJSON{Token: t})
	}
}
