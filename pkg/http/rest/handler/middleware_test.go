package handler_test

import (
	"context"
	"net/http"
	"time"

	"github.com/ifreddyrondon/capture/pkg/domain"
	"github.com/ifreddyrondon/capture/pkg/http/rest/middleware"
	"gopkg.in/src-d/go-kallax.v1"
)

func f2P(v float64) *float64 {
	return &v
}

func s2t(date string) time.Time {
	v, _ := time.Parse(time.RFC3339, date)
	return v
}

var (
	defaultUser    = &domain.User{Email: "test@example.com", ID: kallax.NewULID()}
	defaultCapture = &domain.Capture{ID: kallax.NewULID()}
	defaultRepo    = &domain.Repository{Name: "test public", Visibility: "public"}
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

func withCaptureMiddle(capt *domain.Capture) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			if capt != nil {
				ctx = context.WithValue(ctx, middleware.CaptureCtxKey, capt)
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}
