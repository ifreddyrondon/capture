package middleware_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/go-chi/chi"
	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/middleware/listing"
	"github.com/ifreddyrondon/bastion/render"

	"github.com/ifreddyrondon/bastion/middleware/listing/filtering"

	bastionMiddleware "github.com/ifreddyrondon/bastion/middleware"
	"github.com/ifreddyrondon/bastion/middleware/listing/sorting"
	"github.com/stretchr/testify/assert"
	"gopkg.in/src-d/go-kallax.v1"

	"github.com/ifreddyrondon/capture/pkg/domain"
	"github.com/ifreddyrondon/capture/pkg/http/rest/middleware"
)

var (
	publicVisibility  = filtering.NewValue("public", "public repos")
	privateVisibility = filtering.NewValue("private", "private repos")
	updatedDESC       = sorting.NewSort("updated_at_desc", "updated_at DESC", "Updated date descending")
	updatedASC        = sorting.NewSort("updated_at_asc", "updated_at ASC", "Updated date ascendant")
	createdDESC       = sorting.NewSort("created_at_desc", "created_at DESC", "Created date descending")
	createdASC        = sorting.NewSort("created_at_asc", "created_at ASC", "Created date ascendant")
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

func setupFilterMiddleware(m func(http.Handler) http.Handler) (*bastion.Bastion, *listing.Listing) {
	var result listing.Listing
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l, err := bastionMiddleware.GetListing(r.Context())
		if err != nil {
			render.JSON.InternalServerError(w, err)
			return
		}
		result = *l
		w.Write([]byte("hi"))
	})

	app := bastion.New()
	app.Route("/", func(r chi.Router) {
		r.Use(m)
		r.Get("/", h)
	})
	return app, &result
}
