package handler

import (
	"fmt"
	"net/http"
	"os"

	"github.com/ifreddyrondon/bastion/binder"
	"github.com/ifreddyrondon/bastion/render"
	"github.com/pkg/errors"

	"github.com/ifreddyrondon/capture/pkg/signup"
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
func SignUp(service signup.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload signup.Payload
		if err := binder.JSON.FromReq(r, &payload); err != nil {
			render.JSON.BadRequest(w, err)
			return
		}

		u, err := service.EnrollUser(payload)
		if err != nil {
			if isConflictErr(err) {
				httpErr := render.HTTPError{
					Status:  http.StatusConflict,
					Error:   http.StatusText(http.StatusConflict),
					Message: fmt.Sprintf("email '%v' already exists", *payload.Email),
				}
				render.JSON.Response(w, http.StatusConflict, httpErr)
				return
			}

			fmt.Fprintln(os.Stderr, err)
			render.JSON.InternalServerError(w, err)
			return
		}

		render.JSON.Created(w, u)
	}
}
