package repository

import (
	"context"
	"errors"

	"github.com/ifreddyrondon/capture/pkg"
)

type ctxKey string

const (
	repoKey ctxKey = "repository"
)

var (
	errMissingCapture    = errors.New("repository not found in context")
	errWrongCaptureValue = errors.New("repository value set incorrectly in context")
)

func withRepo(ctx context.Context, repo *pkg.Repository) context.Context {
	return context.WithValue(ctx, repoKey, repo)
}

// GetFromContext will return the repo assigned to the context,
// or nil if there is any error or there isn't a repository.
func GetFromContext(ctx context.Context) (*pkg.Repository, error) {
	tmp := ctx.Value(repoKey)
	if tmp == nil {
		return nil, errMissingCapture
	}
	repo, ok := tmp.(*pkg.Repository)
	if !ok {
		return nil, errWrongCaptureValue
	}
	return repo, nil
}
