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
	// CaptureCtxKey is the context.Context key to store the Capture for a request.
	CaptureCtxKey = &contextKey{"Capture"}
)
var (
	errMissingCtxCapture = errors.New("capture not found in context")
	errWrongCaptureValue = errors.New("capture value set incorrectly in context")
	errMissingCapture    = errors.New("not found capture")
	errInvalidCaptureID  = errors.New("invalid capture id")
)

func withCapture(ctx context.Context, capture *domain.Capture) context.Context {
	return context.WithValue(ctx, CaptureCtxKey, capture)
}

// GetCapture returns the capture assigned to the context, or error if there
// is any error or there isn't a repo.
func GetCapture(ctx context.Context) (*domain.Capture, error) {
	tmp := ctx.Value(CaptureCtxKey)
	if tmp == nil {
		return nil, errMissingCtxCapture
	}
	capt, ok := tmp.(*domain.Capture)
	if !ok {
		return nil, errWrongCaptureValue
	}
	return capt, nil
}

func CaptureCtx(service getting.CaptureService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			captureID := chi.URLParam(r, "captureId")
			repo, err := GetRepo(r.Context())
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				render.JSON.InternalServerError(w, err)
				return
			}

			id, err := kallax.NewULIDFromText(captureID)
			if err != nil {
				render.JSON.BadRequest(w, errInvalidCaptureID)
				return
			}

			capt, err := service.Get(id, repo)
			if err != nil {
				if isNotFound(err) {
					render.JSON.NotFound(w, errMissingCapture)
					return
				}
				fmt.Fprintln(os.Stderr, err)
				render.JSON.InternalServerError(w, err)
				return
			}

			ctx := withCapture(r.Context(), capt)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}
