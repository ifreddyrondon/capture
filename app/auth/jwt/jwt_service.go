package jwt

import (
	"context"
	"errors"
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
	// SubjectCtxKey key to get the subject from requets
	SubjectCtxKey = contextKey("jwt_subject_ctx_key")
)

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

// Service service to managed JWT
type Service struct {
	ExpirationDelta time.Duration
	render          render.Render
	signingKey      []byte
}

// NewService is a helper constructor to create a new service with signing key.
func NewService(signingKey []byte, delta time.Duration, render render.Render) *Service {
	return &Service{
		ExpirationDelta: delta,
		signingKey:      signingKey,
		render:          render,
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

// IsAuthorized validates if a request contains a valid JWT
func (s *Service) IsAuthorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// keyFunc will receive the parsed token and should return the key for validating
		fun := func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrSigningMethod
			}
			return s.signingKey, nil
		}

		token, err := request.ParseFromRequest(
			r, request.OAuth2Extractor,
			fun,
			request.WithClaims(&Claims{}))

		if err != nil || !token.Valid {
			httpErr := json.HTTPError{
				Status:  http.StatusForbidden,
				Errors:  http.StatusText(http.StatusForbidden),
				Message: ErrUserNotAllowed.Error(),
			}
			_ = s.render(w).Response(http.StatusForbidden, httpErr)
			return
		}

		claims, ok := token.Claims.(*Claims)
		if !ok {
			_ = s.render(w).InternalServerError(err)
			return
		}

		ctx := context.WithValue(r.Context(), SubjectCtxKey, claims.Subject)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
