package jwt

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
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
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, NewClaims(userID, s.ExpirationDelta))

	tokenString, err := token.SignedString(s.signingKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
