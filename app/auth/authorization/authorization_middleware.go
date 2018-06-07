package authorization

import (
	"net/http"

	"github.com/ifreddyrondon/bastion/render"
	"github.com/ifreddyrondon/bastion/render/json"
)

// Authorization is a middleware to validate
// if the request can access the resource using a validation strategy
type Authorization struct {
	strategy   Strategy
	render     render.Render
	ctxManager *ContextManager
}

// NewAuthorization returns a new instance of Authorization middleware
func NewAuthorization(strategy Strategy, render render.Render) *Authorization {
	return &Authorization{
		strategy:   strategy,
		render:     render,
		ctxManager: NewContextManager(),
	}
}

// IsAuthorizedREQ validates if a request can access the resource
func (a *Authorization) IsAuthorizedREQ(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		subjectID, err := a.strategy.IsAuthorizedREQ(r)
		if err != nil {
			if a.strategy.IsNotAuthorizedErr(err) {
				httpErr := json.HTTPError{
					Status:  http.StatusForbidden,
					Errors:  http.StatusText(http.StatusForbidden),
					Message: err.Error(),
				}
				_ = a.render(w).Response(http.StatusForbidden, httpErr)
				return
			}
			_ = a.render(w).InternalServerError(err)
			return
		}
		ctx := a.ctxManager.WithSubjectID(r.Context(), subjectID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
