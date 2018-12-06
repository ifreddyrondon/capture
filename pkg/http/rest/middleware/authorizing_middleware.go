package middleware

import (
	"context"
	"net/http"

	"github.com/ifreddyrondon/bastion/render"
	"github.com/ifreddyrondon/capture/pkg"
	"github.com/ifreddyrondon/capture/pkg/authorizing"
	"github.com/pkg/errors"
)

type authorizationErr interface {
	// IsNotAllowed returns true when the req is not allowed.
	IsNotAuthorized() bool
}

func isNotAuthorized(err error) bool {
	if e, ok := errors.Cause(err).(authorizationErr); ok {
		return e.IsNotAuthorized()
	}
	return false
}

type notFoundErr interface {
	// NotFound returns true when a resource is not found.
	NotFound() bool
}

func isNotFound(err error) bool {
	if e, ok := errors.Cause(err).(notFoundErr); ok {
		return e.NotFound()
	}
	return false
}

type ctxKey string

const userKey ctxKey = "user"

var (
	errMissingUser    = errors.New("user not found in context")
	errWrongUserValue = errors.New("user value set incorrectly in context")
)

// WithUser will return a new context with the user value added to it.
func WithUser(ctx context.Context, user *pkg.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

// GetFromContext returns the user assigned to the context, or error if there
// is any error or there isn't a user.
func GetUser(ctx context.Context) (*pkg.User, error) {
	tmp := ctx.Value(userKey)
	if tmp == nil {
		return nil, errMissingUser
	}
	user, ok := tmp.(*pkg.User)
	if !ok {
		return nil, errWrongUserValue
	}
	return user, nil
}

func AuthorizeReq(service authorizing.Service) func(next http.Handler) http.Handler {
	json := render.NewJSON()
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			u, err := service.AuthorizeRequest(r)
			if err != nil {
				if isNotAuthorized(err) {
					httpErr := render.HTTPError{
						Status:  http.StatusForbidden,
						Error:   http.StatusText(http.StatusForbidden),
						Message: "you don’t have permission to access this resource",
					}
					json.Response(w, http.StatusForbidden, httpErr)
					return
				}
				if isNotFound(err) {
					json.NotFound(w, err)
					return
				}
				json.InternalServerError(w, err)
				return
			}

			ctx := WithUser(r.Context(), u)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}
