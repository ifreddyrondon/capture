package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/ifreddyrondon/bastion/render"
	"github.com/ifreddyrondon/capture/pkg/authorizing"
	"github.com/ifreddyrondon/capture/pkg/domain"
	"github.com/pkg/errors"
)

var (
	errInvalidUserID  = errors.New("invalid user id format")
	errMissingUser    = errors.New("user not found in context")
	errWrongUserValue = errors.New("user value set incorrectly in context")
)
var (
	// RepoCtxKey is the context.Context key to store the Repo for a request.
	UserCtxKey = &contextKey{"User"}
)

func withUser(ctx context.Context, user *domain.User) context.Context {
	return context.WithValue(ctx, UserCtxKey, user)
}

// GetUser returns the user assigned to the context, or error if there
// is any error or there isn't a user.
func GetUser(ctx context.Context) (*domain.User, error) {
	tmp := ctx.Value(UserCtxKey)
	if tmp == nil {
		return nil, errMissingUser
	}
	user, ok := tmp.(*domain.User)
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
				if isInvalidErr(err) {
					json.BadRequest(w, errInvalidUserID)
					return
				}
				if isNotAuthorized(err) {
					httpErr := render.HTTPError{
						Status:  http.StatusUnauthorized,
						Error:   http.StatusText(http.StatusUnauthorized),
						Message: "authorization required, access is denied due to invalid credentials",
					}
					json.Response(w, http.StatusUnauthorized, httpErr)
					return
				}
				if isNotFound(err) {
					json.NotFound(w, err)
					return
				}
				fmt.Fprintln(os.Stderr, err)
				json.InternalServerError(w, err)
				return
			}

			ctx := withUser(r.Context(), u)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}
