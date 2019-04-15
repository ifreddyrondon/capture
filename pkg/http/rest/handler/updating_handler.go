package handler

import (
	"fmt"
	"net/http"
	"os"

	"github.com/ifreddyrondon/bastion/render"

	"github.com/ifreddyrondon/capture/pkg/http/rest/middleware"
	"github.com/ifreddyrondon/capture/pkg/updating"
)

// UpdatingCapture returns a configured http.Handler with updating capture resources.
func UpdatingCapture(service updating.CaptureService) http.HandlerFunc {
	renderJSON := render.NewJSON()

	return func(w http.ResponseWriter, r *http.Request) {
		capt, err := middleware.GetCapture(r.Context())
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			renderJSON.InternalServerError(w, err)
			return
		}

		var data updating.Capture
		err = updating.CaptureValidator.Decode(r, &data)
		if err != nil {
			renderJSON.BadRequest(w, err)
			return
		}

		err = service.Update(data, capt)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			renderJSON.InternalServerError(w, err)
			return
		}

		renderJSON.Send(w, capt)
	}
}
