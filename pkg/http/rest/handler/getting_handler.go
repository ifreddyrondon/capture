package handler

import (
	"fmt"
	"net/http"
	"os"

	"github.com/ifreddyrondon/bastion/render"

	"github.com/ifreddyrondon/capture/pkg/http/rest/middleware"
)

// GettingRepo returns a configured http.Handler with getting repo resources.
func GettingRepo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		repo, err := middleware.GetRepo(r.Context())
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			render.JSON.InternalServerError(w, err)
			return
		}

		render.JSON.Send(w, repo)
	}
}

// GettingCapture returns a configured http.Handler with getting capture resources.
func GettingCapture() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		capt, err := middleware.GetCapture(r.Context())
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			render.JSON.InternalServerError(w, err)
			return
		}

		render.JSON.Send(w, capt)
	}
}
