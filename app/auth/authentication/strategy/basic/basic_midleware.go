package basic

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ifreddyrondon/capture/app/user"

	"github.com/ifreddyrondon/bastion/render"
	bastionJSON "github.com/ifreddyrondon/bastion/render/json"
)

var (
	errInvalidCredentials = errors.New("invalid email or password")
	errInvalidPassword    = errors.New("invalid password")
)

// Strategy is a basic authentication method that uses email and password to authenticate
type Strategy struct {
	render     render.Render
	service    user.GetterService
	ctxManager *user.ContextManager
}

// NewStrategy returns a new instance of Strategy
func NewStrategy(render render.Render, service user.GetterService) *Strategy {
	return &Strategy{render: render, service: service, ctxManager: user.NewContextManager()}
}

// Authenticate for basic (username/password) authentication.
func (s *Strategy) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var cre Crendentials
		if err := json.NewDecoder(r.Body).Decode(&cre); err != nil {
			_ = s.render(w).BadRequest(err)
			return
		}

		u, err := s.validate(&cre)
		if err != nil {
			if err == errInvalidPassword || err == user.ErrNotFound {
				httpErr := bastionJSON.HTTPError{
					Status:  http.StatusUnauthorized,
					Errors:  http.StatusText(http.StatusUnauthorized),
					Message: errInvalidCredentials.Error(),
				}
				_ = s.render(w).Response(http.StatusUnauthorized, httpErr)
				return
			}

			_ = s.render(w).InternalServerError(err)
			return
		}
		ctx := s.ctxManager.WithUser(r.Context(), u)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *Strategy) validate(cre *Crendentials) (*user.User, error) {
	u, err := s.service.GetByEmail(cre.Email)
	if err != nil {
		return nil, err
	}
	if !u.CheckPassword(cre.Password) {
		return nil, errInvalidPassword
	}

	return u, nil
}
