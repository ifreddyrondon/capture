package rest

import (
	"fmt"
	"net/http"
	"os"

	"github.com/ifreddyrondon/bastion/render"
	"github.com/ifreddyrondon/capture/pkg/adding"
	"github.com/ifreddyrondon/capture/pkg/http/rest/middleware"
)

// Adding returns a configured http.Handler with adding resources.
func Adding(service adding.Service) http.HandlerFunc {
	renderJSON := render.NewJSON()

	return func(w http.ResponseWriter, r *http.Request) {
		var payl adding.Capture
		err := adding.CaptureValidator.Decode(r, &payl)
		if err != nil {
			renderJSON.BadRequest(w, err)
			return
		}

		repo, err := middleware.GetRepo(r.Context())
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			renderJSON.InternalServerError(w, err)
			return
		}

		capt, err := service.AddCapture(repo, payl)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			renderJSON.InternalServerError(w, err)
			return
		}

		renderJSON.Created(w, capt)
	}
}
