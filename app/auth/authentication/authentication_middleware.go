package authentication

import (
	"encoding/json"
	"net/http"

	"github.com/ifreddyrondon/bastion/render"
	bastionJSON "github.com/ifreddyrondon/bastion/render/json"
	"github.com/ifreddyrondon/capture/app/user"
)

// Authentication is a middleware to validate request credentials.
type Authentication struct {
	strategy   Strategy
	render     render.Render
	ctxManager *user.ContextManager
}

// NewAuthentication returns a new instance of Authentication middleware
func NewAuthentication(strategy Strategy, render render.Render) *Authentication {
	return &Authentication{
		strategy:   strategy,
		render:     render,
		ctxManager: user.NewContextManager(),
	}
}

// Authenticate validate if an user is authorized to continue or 401.
func (a *Authentication) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload json.RawMessage
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			_ = a.render(w).BadRequest(err)
			return
		}

		u, err := a.strategy.Validate(payload)
		if err != nil {
			if a.strategy.IsErrCredentials(err) {
				httpErr := bastionJSON.HTTPError{
					Status:  http.StatusUnauthorized,
					Errors:  http.StatusText(http.StatusUnauthorized),
					Message: err.Error(),
				}
				_ = a.render(w).Response(http.StatusUnauthorized, httpErr)
				return
			}
			if a.strategy.IsErrDecoding(err) {
				_ = a.render(w).BadRequest(err)
				return
			}
			_ = a.render(w).InternalServerError(err)
			return
		}
		ctx := a.ctxManager.WithUser(r.Context(), u)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
