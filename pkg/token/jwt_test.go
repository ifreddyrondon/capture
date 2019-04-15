package token_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/ifreddyrondon/capture/pkg"

	"github.com/dgrijalva/jwt-go"

	"github.com/ifreddyrondon/capture/pkg/token"
)

func TestJWTClaims(t *testing.T) {
	t.Parallel()

	expected := time.Date(1989, time.Month(12), 26, 6, 1, 0, 0, time.UTC)
	mockClock := pkg.NewMockClock(expected)
	userID := "123"

	tt := []struct {
		name            string
		expirationDelta int64
		expect          token.JWTClaims
	}{
		{
			"claims with default delta",
			0,
			token.JWTClaims{
				StandardClaims: jwt.StandardClaims{
					Subject:  userID,
					IssuedAt: expected.Unix(),
					ExpiresAt: func() int64 {
						n := expected
						return n.Add(token.DefaultJWTExpirationDelta).Unix()
					}(),
				},
			},
		},
		{
			"claims with 30 min delta",
			1800, // 30min
			token.JWTClaims{
				StandardClaims: jwt.StandardClaims{
					Subject:  userID,
					IssuedAt: expected.Unix(),
					ExpiresAt: func() int64 {
						n := expected
						return n.Add(time.Minute * 30).Unix()
					}(),
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var c token.JWTClaims
			// mock clock
			token.SetClockInstance(&c, mockClock)
			c.Subject = userID
			c.IssueIt()
			c.SetExpirationDate(0)

			assert.Equal(t, userID, c.Subject)
		})
	}
}

func TestInvalidClaims(t *testing.T) {
	t.Parallel()

	userID := "123"
	tt := []struct {
		name   string
		claims token.JWTClaims
		err    string
	}{
		{
			"expired",
			token.JWTClaims{
				StandardClaims: jwt.StandardClaims{
					Subject:   userID,
					IssuedAt:  time.Now().Unix(),
					ExpiresAt: time.Date(1989, time.Month(12), 26, 6, 1, 0, 0, time.UTC).Unix(),
				},
			},
			"token is expired by",
		},
		{
			"issued after valid",
			token.JWTClaims{
				StandardClaims: jwt.StandardClaims{
					Subject:   userID,
					IssuedAt:  time.Now().Add(time.Minute * 10).Unix(),
					ExpiresAt: time.Now().Add(time.Minute * 10).Unix(),
				},
			},
			"Token used before issued",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.claims.Valid()
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), tc.err)
		})
	}
}

func TestNewClaims(t *testing.T) {
	t.Parallel()

	c := token.NewJWTClaims("123", 0)
	err := c.Valid()
	assert.Nil(t, err)
}
