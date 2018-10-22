package jwt

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/ifreddyrondon/capture/features"
)

// DefaultJWTExpirationDelta is the delta added to time.Now() when a claims is created
const DefaultJWTExpirationDelta time.Duration = time.Hour

// Claims Section of JWT. Referenced at https://tools.ietf.org/html/rfc7519#section-4.1
type Claims struct {
	jwt.StandardClaims
	clock *features.Clock
}

// IssueIt marks the claims with iat (IssuedAt).
func (c *Claims) IssueIt() {
	c.IssuedAt = c.clock.Now().Unix()
}

// SetExpirationDate marks the claims with exp (ExpiresAt).
// If expirationDelta is 0 (default) then it takes DefaultJWTExpirationDelta (30min).
func (c *Claims) SetExpirationDate(delta time.Duration) {
	if delta == 0 {
		delta = DefaultJWTExpirationDelta
	}
	c.ExpiresAt = c.clock.Now().Add(delta).Unix()
}

// Valid is a wrapper over jwt.StandardClaims valid
func (c *Claims) Valid() error {
	return c.StandardClaims.Valid()
}

// NewClaims returns a new filled jwt claims.
// Subject filled with user id.
// IssuedAt filled with time.Now()
// ExpiresAt filled with a future exp date provided by arg expirationDelta or DefaultJWTExpirationDelta
func NewClaims(userID string, delta time.Duration) jwt.Claims {
	c := new(Claims)
	c.Subject = userID
	c.IssueIt()
	c.SetExpirationDate(delta)

	return c
}
