package authentication

import (
	"net/http"

	"github.com/ifreddyrondon/bastion/render"
	"github.com/ifreddyrondon/capture/app/user"
)

// Authentication is a middleware to validate request credentials.
type Authentication struct {
	strategy   Strategy
	render     render.APIRenderer
	ctxManager *user.ContextManager
}

// NewAuthentication returns a new instance of Authentication middleware
func NewAuthentication(strategy Strategy) *Authentication {
	return &Authentication{
		strategy:   strategy,
		render:     render.NewJSON(),
		ctxManager: user.NewContextManager(),
	}
}

// Authenticate validate if an user is authorized to continue or 401.
func (a *Authentication) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, err := a.strategy.Validate(r)
		if err != nil {
			if a.strategy.IsErrCredentials(err) {
				httpErr := render.HTTPError{
					Status:  http.StatusUnauthorized,
					Error:   http.StatusText(http.StatusUnauthorized),
					Message: err.Error(),
				}
				a.render.Response(w, http.StatusUnauthorized, httpErr)
				return
			}
			if a.strategy.IsErrDecoding(err) {
				a.render.BadRequest(w, err)
				return
			}
			a.render.InternalServerError(w, err)
			return
		}
		ctx := a.ctxManager.WithUser(r.Context(), u)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
