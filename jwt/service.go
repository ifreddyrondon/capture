package jwt

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ifreddyrondon/bastion/render"
	"github.com/ifreddyrondon/bastion/render/json"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
)

var (
	// ErrSigningMethod expected error when signing method is not the same when the token was signed
	ErrSigningMethod = errors.New("unexpected signing method")
	// ErrUserNotAllowed expected error when a user tries to access a resource that is not his
	ErrUserNotAllowed = errors.New("you don’t have permission to access this resource")
	// SigningMethod method used or to be used in jwt
	SigningMethod = jwt.SigningMethodHS256
)

// Service service to managed JWT
type Service struct {
	ExpirationDelta time.Duration
	UserIDKey       fmt.Stringer
	render.Render
	signingKey []byte
}

// NewService is a helper constructor to create a new service with signing key.
func NewService(signingKey []byte, delta time.Duration) *Service {
	return &Service{
		ExpirationDelta: delta,
		signingKey:      signingKey,
	}
}

// GenerateToken creates a new JWT
func (s *Service) GenerateToken(userID string) (string, error) {
	token := jwt.NewWithClaims(SigningMethod, NewClaims(userID, s.ExpirationDelta))

	tokenString, err := token.SignedString(s.signingKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Authorization validates if a request contains a valid JWT
func (s *Service) Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := request.ParseFromRequest(r, request.OAuth2Extractor, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrSigningMethod
			}

			return s.signingKey, nil
		})

		if err != nil || !token.Valid {
			httpErr := json.HTTPError{
				Status:  http.StatusForbidden,
				Errors:  http.StatusText(http.StatusForbidden),
				Message: ErrUserNotAllowed.Error(),
			}
			_ = s.Render(w).Response(http.StatusForbidden, httpErr)
			return
		}

		claims, ok := token.Claims.(*Claims)
		if !ok {
			_ = s.Render(w).InternalServerError(err)
			return
		}

		ctx := context.WithValue(r.Context(), s.UserIDKey, claims.Subject)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
