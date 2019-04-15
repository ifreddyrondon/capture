package handler

import (
	"fmt"
	"net/http"
	"os"

	"github.com/ifreddyrondon/bastion/render"

	"github.com/ifreddyrondon/capture/pkg/creating"
	"github.com/ifreddyrondon/capture/pkg/http/rest/middleware"
)

// Creating returns a configured http.Handler with creating resources.
func Creating(service creating.Service) http.HandlerFunc {
	renderJSON := render.NewJSON()

	return func(w http.ResponseWriter, r *http.Request) {
		var payl creating.Payload
		err := creating.Validator.Decode(r, &payl)
		if err != nil {
			renderJSON.BadRequest(w, err)
			return
		}

		u, err := middleware.GetUser(r.Context())
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			renderJSON.InternalServerError(w, err)
			return
		}

		repo, err := service.CreateRepo(u, payl)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			renderJSON.InternalServerError(w, err)
			return
		}

		renderJSON.Created(w, repo)
	}
}
