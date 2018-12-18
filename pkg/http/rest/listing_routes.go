package rest

import (
	"fmt"
	"net/http"
	"os"

	"github.com/ifreddyrondon/bastion/middleware"
	"github.com/ifreddyrondon/bastion/render"
	auth "github.com/ifreddyrondon/capture/pkg/http/rest/middleware"
	"github.com/ifreddyrondon/capture/pkg/listing"
)

// ListingUserRepos returns a configured http.Handler with user repos resources to get user's repos.
func ListingUserRepos(service listing.Service) http.HandlerFunc {
	renderJSON := render.NewJSON()
	return func(w http.ResponseWriter, r *http.Request) {
		l, err := middleware.GetListing(r.Context())
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			renderJSON.InternalServerError(w, err)
			return
		}

		u, err := auth.GetUser(r.Context())
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			renderJSON.InternalServerError(w, err)
			return
		}

		res, err := service.GetUserRepos(u, l)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			renderJSON.InternalServerError(w, err)
			return
		}

		renderJSON.Send(w, res)
	}
}

// ListingPublicRepos returns a configured http.Handler with repos resources to get public repos.
func ListingPublicRepos(service listing.Service) http.HandlerFunc {
	renderJSON := render.NewJSON()
	return func(w http.ResponseWriter, r *http.Request) {
		l, err := middleware.GetListing(r.Context())
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			renderJSON.InternalServerError(w, err)
			return
		}

		res, err := service.GetPublicRepos(l)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			renderJSON.InternalServerError(w, err)
			return
		}

		renderJSON.Send(w, res)
	}
}
