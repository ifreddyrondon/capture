package rest

import (
	"fmt"
	"net/http"
	"os"

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
func SignUp(service signup.Service) http.HandlerFunc {
	renderJSON := render.NewJSON()
	return func(w http.ResponseWriter, r *http.Request) {
		var payl signup.Payload
		err := signup.Validator.Decode(r, &payl)
		if err != nil {
			renderJSON.BadRequest(w, err)
			return
		}

		u, err := service.EnrollUser(payl)
		if err != nil {
			if isConflictErr(err) {
				httpErr := render.HTTPError{
					Status:  http.StatusConflict,
					Error:   http.StatusText(http.StatusConflict),
					Message: fmt.Sprintf("email '%v' already exists", *payl.Email),
				}
				renderJSON.Response(w, http.StatusConflict, httpErr)
				return
			}

			fmt.Fprintln(os.Stderr, err)
			renderJSON.InternalServerError(w, err)
			return
		}

		renderJSON.Created(w, u)
	}
}
