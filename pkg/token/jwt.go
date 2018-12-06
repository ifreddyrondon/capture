package token

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/ifreddyrondon/capture/pkg"
)

// DefaultJWTExpirationDelta is the delta added to time.Now() when a jwt claims is created
const DefaultJWTExpirationDelta = time.Hour

// JWTClaims Section of JWT. Referenced at https://tools.ietf.org/html/rfc7519#section-4.1
type JWTClaims struct {
	jwt.StandardClaims
	clock *pkg.Clock
}

// IssueIt marks the claims with iat (IssuedAt).
func (c *JWTClaims) IssueIt() {
	c.IssuedAt = c.clock.Now().Unix()
}

// SetExpirationDate marks the claims with exp (ExpiresAt).
// If expirationDelta is 0 (default) then it takes DefaultJWTExpirationDelta (30min).
func (c *JWTClaims) SetExpirationDate(delta time.Duration) {
	if delta == 0 {
		delta = DefaultJWTExpirationDelta
	}
	c.ExpiresAt = c.clock.Now().Add(delta).Unix()
}

// Valid is a wrapper over jwt.StandardClaims valid
func (c *JWTClaims) Valid() error {
	return c.StandardClaims.Valid()
}

// NewJWTClaims returns a new filled jwt claims.
// Subject filled with user id.
// IssuedAt filled with time.Now()
// ExpiresAt filled with a future exp date provided by arg expirationDelta or DefaultJWTExpirationDelta
func NewJWTClaims(userID string, delta time.Duration) jwt.Claims {
	c := new(JWTClaims)
	c.Subject = userID
	c.IssueIt()
	c.SetExpirationDate(delta)

	return c
}
