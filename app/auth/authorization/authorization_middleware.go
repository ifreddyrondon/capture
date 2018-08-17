package authorization

import (
	"net/http"

	"github.com/ifreddyrondon/bastion/render"
)

// Authorization is a middleware to validate
// if the request can access the resource using a validation strategy
type Authorization struct {
	strategy   Strategy
	render     render.APIRenderer
	ctxManager *ContextManager
}

// NewAuthorization returns a new instance of Authorization middleware
func NewAuthorization(strategy Strategy) *Authorization {
	return &Authorization{
		strategy:   strategy,
		render:     render.NewJSON(),
		ctxManager: NewContextManager(),
	}
}

// IsAuthorizedREQ validates if a request can access the resource
func (a *Authorization) IsAuthorizedREQ(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		subjectID, err := a.strategy.IsAuthorizedREQ(r)
		if err != nil {
			if a.strategy.IsNotAuthorizedErr(err) {
				httpErr := render.HTTPError{
					Status:  http.StatusForbidden,
					Error:   http.StatusText(http.StatusForbidden),
					Message: err.Error(),
				}
				a.render.Response(w, http.StatusForbidden, httpErr)
				return
			}
			a.render.InternalServerError(w, err)
			return
		}
		ctx := a.ctxManager.WithSubjectID(r.Context(), subjectID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
