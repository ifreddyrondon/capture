package user

import (
	"net/http"

	"gopkg.in/src-d/go-kallax.v1"

	"github.com/ifreddyrondon/bastion/render"
	"github.com/ifreddyrondon/capture/app/auth/authorization"
)

// Middleware are helper methods to set user information into a request context.
type Middleware struct {
	authorizationCtxManager *authorization.ContextManager
	userCtxManager          *ContextManager
	service                 GetterService
	render                  render.APIRenderer
}

// NewMiddleware returns a new instance of Middleware
func NewMiddleware(service GetterService) *Middleware {
	return &Middleware{
		authorizationCtxManager: authorization.NewContextManager(),
		service:                 service,
		render:                  render.NewJSON(),
	}
}

// LoggedUser save the authenticated user in a request context.
func (m *Middleware) LoggedUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		subjectID := m.authorizationCtxManager.Get(r.Context())
		userID, err := kallax.NewULIDFromText(subjectID)
		if err != nil {
			m.render.BadRequest(w, err)
			return
		}
		u, err := m.service.GetByID(userID)
		if err != nil {
			if err == ErrNotFound {
				m.render.NotFound(w, err)
				return
			}
			m.render.InternalServerError(w, err)
			return
		}
		ctx := m.userCtxManager.WithUser(r.Context(), u)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
