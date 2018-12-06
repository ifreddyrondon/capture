package token

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/pkg/errors"
)

var (
	errSigningMethod = errors.New("unexpected signing method")
	errInvalidClaims = errors.New("error invalid jwt claims")
)

type notAllowedErr string

func (i notAllowedErr) Error() string         { return fmt.Sprintf(string(i)) }
func (i notAllowedErr) IsNotAuthorized() bool { return true }

type JwtService struct {
	expirationDelta time.Duration
	signingKey      []byte
}

// NewJWTService is a helper constructor to create a new service with signing key.
func NewJWTService(signingKey string, delta time.Duration) *JwtService {
	return &JwtService{
		expirationDelta: delta,
		signingKey:      []byte(signingKey),
	}
}

// GenerateToken creates a new JWT
func (s *JwtService) GenerateToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, NewJWTClaims(userID, s.expirationDelta))

	tokenString, err := token.SignedString(s.signingKey)
	if err != nil {
		return "", errors.Wrap(err, "GenerateToken")
	}

	return tokenString, nil
}

func (s *JwtService) IsRequestAuthorized(r *http.Request) (string, error) {
	token, err := request.ParseFromRequest(
		r, request.OAuth2Extractor, s.validateMethod, request.WithClaims(&JWTClaims{}))

	if err != nil || !token.Valid {
		return "", errors.WithStack(notAllowedErr(err.Error()))
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return "", errInvalidClaims
	}

	return claims.Subject, nil
}

// validateMethod will receive the parsed token and should return the key for validating
func (s *JwtService) validateMethod(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, errSigningMethod
	}
	return s.signingKey, nil
}
