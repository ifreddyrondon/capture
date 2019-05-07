package handler

import (
	"fmt"
	"net/http"
	"os"

	"github.com/ifreddyrondon/bastion/binder"
	"github.com/ifreddyrondon/bastion/render"

	"github.com/ifreddyrondon/capture/pkg/creating"
	"github.com/ifreddyrondon/capture/pkg/http/rest/middleware"
)

// Creating returns a configured http.Handler with creating resources.
func Creating(service creating.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload creating.Payload
		if err := binder.JSON.FromReq(r, &payload); err != nil {
			render.JSON.BadRequest(w, err)
			return
		}

		u, err := middleware.GetUser(r.Context())
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			render.JSON.InternalServerError(w, err)
			return
		}

		repo, err := service.CreateRepo(u, payload)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			render.JSON.InternalServerError(w, err)
			return
		}

		render.JSON.Created(w, repo)
	}
}
