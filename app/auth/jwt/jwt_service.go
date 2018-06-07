package jwt

import (
	"errors"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
)

var (
	errSigningMethod  = errors.New("unexpected signing method")
	errUserNotAllowed = errors.New("you don’t have permission to access this resource")
	errInvalidClaims  = errors.New("error invalid jwt claims")
	// SigningMethod method used or to be used in jwt
	SigningMethod = jwt.SigningMethodHS256
)

// Service service to managed JWT
type Service struct {
	ExpirationDelta time.Duration
	signingKey      []byte
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

// IsAuthorizedREQ validates if a request contains a valid JWT
func (s *Service) IsAuthorizedREQ(r *http.Request) (string, error) {
	token, err := request.ParseFromRequest(
		r, request.OAuth2Extractor, s.validateMethod, request.WithClaims(&Claims{}))

	if err != nil || !token.Valid {
		return "", errUserNotAllowed
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return "", errInvalidClaims
	}

	return claims.Subject, nil
}

// IsNotAuthorizedErr check if an error is for invalid jwt.
func (s *Service) IsNotAuthorizedErr(err error) bool {
	return err == errUserNotAllowed
}

// validateMethod will receive the parsed token and should return the key for validating
func (s *Service) validateMethod(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, errSigningMethod
	}
	return s.signingKey, nil
}
