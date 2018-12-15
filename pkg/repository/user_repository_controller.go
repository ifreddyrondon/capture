package repository

import (
	"errors"
	"net/http"

	"github.com/ifreddyrondon/bastion/render"

	"github.com/ifreddyrondon/bastion/middleware"
	auth "github.com/ifreddyrondon/capture/pkg/http/rest/middleware"
	"github.com/ifreddyrondon/capture/pkg/repository/encoder"
)

var (
	// ErrorInvalidRepoID expected error when repository id param is invalid
	ErrorInvalidRepoID = errors.New("invalid repository id")
	// ErrorNotFound expected error when repository is missing
	ErrorNotFound = errors.New("not found repository")
	// ErrorNotFound expected error when repository is missing
	ErrorNotAuthorized = errors.New("not authorized to see this repository")
)

// UserRoutes returns a configured http.Handler with user repositories resources.
func ListingOwnRepos(service Service) http.HandlerFunc {
	renderJSON := render.NewJSON()
	return func(w http.ResponseWriter, r *http.Request) {
		listing, err := middleware.GetListing(r.Context())
		if err != nil {
			renderJSON.InternalServerError(w, err)
			return
		}

		u, err := auth.GetUser(r.Context())
		if err != nil {
			renderJSON.InternalServerError(w, err)
			return
		}

		repos, err := service.GetUserRepositories(u, listing)
		if err != nil {
			renderJSON.InternalServerError(w, err)
			return
		}

		res := encoder.ListRepositoryResponse{Listing: listing, Results: repos}
		renderJSON.Send(w, res)
	}
}
