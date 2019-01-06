package rest_test

import (
	"context"
	"net/http"

	"github.com/ifreddyrondon/capture/pkg/domain"
	"github.com/ifreddyrondon/capture/pkg/http/rest/middleware"
	"gopkg.in/src-d/go-kallax.v1"
)

var defaultCapture = &domain.Capture{ID: kallax.NewULID()}

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
