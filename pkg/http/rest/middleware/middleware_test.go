package middleware_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/ifreddyrondon/capture/pkg/domain"
	"github.com/ifreddyrondon/capture/pkg/http/rest/middleware"
	"github.com/stretchr/testify/assert"
	kallax "gopkg.in/src-d/go-kallax.v1"
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

var tempUser = domain.User{Email: "test@example.com", ID: kallax.NewULID()}

func authenticatedMiddle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), middleware.UserCtxKey, &tempUser)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func notUserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

var tempRepo = domain.Repository{Name: "test", ID: kallax.NewULID()}

func getRepoOKMiddle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), middleware.RepoCtxKey, &tempRepo)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func notRepoMiddle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}
