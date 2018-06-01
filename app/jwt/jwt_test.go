package jwt_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/dgrijalva/jwt-go"
	captureJWT "github.com/ifreddyrondon/capture/app/jwt"
	"github.com/ifreddyrondon/capture/app/timestamp"
)

func TestClaims(t *testing.T) {
	t.Parallel()

	expected := time.Date(1989, time.Month(12), 26, 6, 1, 0, 0, time.UTC)
	mockClock := timestamp.NewMockClock(expected)
	userID := "123"

	tt := []struct {
		name            string
		expirationDelta int64
		expect          captureJWT.Claims
	}{
		{
			"claims with default delta",
			0,
			captureJWT.Claims{
				StandardClaims: jwt.StandardClaims{
					Subject:  userID,
					IssuedAt: expected.Unix(),
					ExpiresAt: func() int64 {
						n := expected
						return n.Add(captureJWT.DefaultJWTExpirationDelta).Unix()
					}(),
				},
			},
		},
		{
			"claims with 30 min delta",
			1800, // 30min
			captureJWT.Claims{
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
			var c captureJWT.Claims
			// mock clock
			captureJWT.SetClockInstance(&c, mockClock)
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
		claims captureJWT.Claims
		err    string
	}{
		{
			"expired",
			captureJWT.Claims{
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
			captureJWT.Claims{
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

	c := captureJWT.NewClaims("123", 0)
	err := c.Valid()
	assert.Nil(t, err)
}
