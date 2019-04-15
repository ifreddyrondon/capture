package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/ifreddyrondon/bastion/render"
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-kallax.v1"

	"github.com/ifreddyrondon/capture/pkg/domain"
	"github.com/ifreddyrondon/capture/pkg/getting"
)

var (
	// RepoCtxKey is the context.Context key to store the Repo for a request.
	RepoCtxKey = &contextKey{"Repository"}
)
var (
	errMissingCtxRepo = errors.New("repo not found in context")
	errWrongRepoValue = errors.New("repo value set incorrectly in context")
	errMissingRepo    = errors.New("not found repository")
	errInvalidRepoID  = errors.New("invalid repository id")
)

func withRepo(ctx context.Context, repo *domain.Repository) context.Context {
	return context.WithValue(ctx, RepoCtxKey, repo)
}

// GetRepo returns the repo assigned to the context, or error if there
// is any error or there isn't a repo.
func GetRepo(ctx context.Context) (*domain.Repository, error) {
	tmp := ctx.Value(RepoCtxKey)
	if tmp == nil {
		return nil, errMissingCtxRepo
	}
	repo, ok := tmp.(*domain.Repository)
	if !ok {
		return nil, errWrongRepoValue
	}
	return repo, nil
}

func RepoCtx(service getting.RepoService) func(next http.Handler) http.Handler {
	json := render.NewJSON()
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			repoID := chi.URLParam(r, "id")
			id, err := kallax.NewULIDFromText(repoID)
			if err != nil {
				json.BadRequest(w, errInvalidRepoID)
				return
			}

			repo, err := service.Get(id)
			if err != nil {
				if isNotFound(err) {
					json.NotFound(w, errMissingRepo)
					return
				}
				fmt.Fprintln(os.Stderr, err)
				json.InternalServerError(w, err)
				return
			}

			ctx := withRepo(r.Context(), repo)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}
