package handler

import (
	"fmt"
	"net/http"
	"os"

	"github.com/ifreddyrondon/bastion/render"

	"github.com/ifreddyrondon/capture/pkg/http/rest/middleware"
	"github.com/ifreddyrondon/capture/pkg/removing"
)

// RemovingCapture returns a configured http.Handler with removing capture resources.
func RemovingCapture(service removing.CaptureService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		capt, err := middleware.GetCapture(r.Context())
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			render.JSON.InternalServerError(w, err)
			return
		}

		err = service.Remove(capt)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			render.JSON.InternalServerError(w, err)
			return
		}

		render.JSON.Send(w, capt)
	}
}
