package middleware_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/ifreddyrondon/capture/pkg/domain"
	"github.com/ifreddyrondon/capture/pkg/http/rest/middleware"
	"github.com/stretchr/testify/assert"
	"gopkg.in/src-d/go-kallax.v1"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "ok")
}

type notFound string

func (u notFound) Error() string  { return fmt.Sprintf(string(u)) }
func (u notFound) NotFound() bool { return true }

type invalidErr string

func (i invalidErr) Error() string   { return fmt.Sprintf(string(i)) }
func (i invalidErr) IsInvalid() bool { return true }

type notAllowedErr string

func (i notAllowedErr) Error() string         { return fmt.Sprintf(string(i)) }
func (i notAllowedErr) IsNotAuthorized() bool { return true }

func TestContextKeyString(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "capture/middleware context value Repository", middleware.RepoCtxKey.String())
}

var (
	defaultUserID = kallax.NewULID()
	defaultUser   = &domain.User{Email: "test@example.com", ID: defaultUserID}
	defaultRepo   = &domain.Repository{Name: "test", ID: kallax.NewULID()}
)

func withUserMiddle(user *domain.User) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			if user != nil {
				ctx = context.WithValue(ctx, middleware.UserCtxKey, user)
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

func withRepoMiddle(repo *domain.Repository) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			if repo != nil {
				ctx = context.WithValue(ctx, middleware.RepoCtxKey, repo)
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}
