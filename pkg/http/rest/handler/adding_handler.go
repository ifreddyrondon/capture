package handler

import (
	"fmt"
	"net/http"
	"os"

	"github.com/ifreddyrondon/bastion/binder"
	"github.com/ifreddyrondon/bastion/render"

	"github.com/ifreddyrondon/capture/pkg/adding"
	"github.com/ifreddyrondon/capture/pkg/http/rest/middleware"
)

// AddingCapture returns a configured http.Handler with adding capture resources.
func AddingCapture(service adding.CaptureService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload adding.Capture
		if err := binder.JSON.FromReq(r, &payload); err != nil {
			render.JSON.BadRequest(w, err)
			return
		}

		repo, err := middleware.GetRepo(r.Context())
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			render.JSON.InternalServerError(w, err)
			return
		}

		capt, err := service.AddCapture(repo, payload)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			render.JSON.InternalServerError(w, err)
			return
		}

		render.JSON.Created(w, capt)
	}
}

// AddingMultiCapture returns a configured http.Handler with adding captures resources.
func AddingMultiCapture(service adding.MultiCaptureService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var multi adding.MultiCapture
		if err := binder.JSON.FromReq(r, &multi); err != nil {
			render.JSON.BadRequest(w, err)
			return
		}

		repo, err := middleware.GetRepo(r.Context())
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			render.JSON.InternalServerError(w, err)
			return
		}

		captures, err := service.AddCaptures(repo, multi)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			render.JSON.InternalServerError(w, err)
			return
		}

		render.JSON.Created(w, captures)
	}
}
