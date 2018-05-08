package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ifreddyrondon/bastion/render"

	bastionJSON "github.com/ifreddyrondon/bastion/render/json"
	"github.com/ifreddyrondon/gocapture/user"
)

// Strategies are Authentication mechanisms to validate users credentials
type Strategies struct {
	render.Render
	Service
	CtxKey fmt.Stringer
}

// LocalStrategy for username/password authentication
func (s *Strategies) LocalStrategy(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var credentials BasicAuthCrendential
		if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
			_ = s.Render(w).BadRequest(err)
			return
		}

		u, err := s.Authenticate(&credentials)
		if err != nil {
			if err == errInvalidPassword || err == user.ErrNotFound {
				httpErr := bastionJSON.HTTPError{
					Status:  http.StatusUnauthorized,
					Errors:  http.StatusText(http.StatusUnauthorized),
					Message: errInvalidCredentials.Error(),
				}
				_ = s.Render(w).Response(http.StatusUnauthorized, httpErr)
				return
			}

			_ = s.Render(w).InternalServerError(err)
			return
		}
		ctx := context.WithValue(r.Context(), s.CtxKey, u)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
