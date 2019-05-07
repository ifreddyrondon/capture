package middleware

import (
	"fmt"
	"net/http"
	"os"

	"github.com/ifreddyrondon/bastion/render"

	"github.com/ifreddyrondon/capture/pkg/authorizing"
	"github.com/ifreddyrondon/capture/pkg/domain"
)

func forbidden(errMsg string) render.HTTPError {
	return render.HTTPError{
		Status:  http.StatusForbidden,
		Error:   http.StatusText(http.StatusForbidden),
		Message: errMsg,
	}
}

func getUserAndRepo(r *http.Request) (*domain.User, *domain.Repository, error) {
	u, err := GetUser(r.Context())
	if err != nil {
		return nil, nil, err
	}

	repo, err := GetRepo(r.Context())
	if err != nil {
		return nil, nil, err
	}

	return u, repo, nil
}

func RepoOwnerOrPublic() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			u, repo, err := getUserAndRepo(r)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				render.JSON.InternalServerError(w, err)
				return
			}
			p := authorizing.RepoPermission(*repo)
			if !p.IsOwnerOrPublic(u.ID) {
				errMsg := fmt.Sprintf("You don't have permission to access repository %s", repo.ID)
				render.JSON.Response(w, http.StatusForbidden, forbidden(errMsg))
				return
			}

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

func RepoOwner() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			u, repo, err := getUserAndRepo(r)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				render.JSON.InternalServerError(w, err)
				return
			}
			p := authorizing.RepoPermission(*repo)
			if !p.IsOwner(u.ID) {
				errMsg := fmt.Sprintf("You don't have permission to access repository %s", repo.ID)
				render.JSON.Response(w, http.StatusForbidden, forbidden(errMsg))
				return
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
