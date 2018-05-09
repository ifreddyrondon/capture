package jwt_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/dgrijalva/jwt-go"
	gocaptureJWT "github.com/ifreddyrondon/gocapture/jwt"
	"github.com/ifreddyrondon/gocapture/timestamp"
)

func TestClaims(t *testing.T) {
	t.Parallel()

	expected := time.Date(1989, time.Month(12), 26, 6, 1, 0, 0, time.UTC)
	mockClock := timestamp.NewMockClock(expected)
	userID := "123"

	tt := []struct {
		name            string
		expirationDelta int64
		expect          gocaptureJWT.Claims
	}{
		{
			"claims with default delta",
			0,
			gocaptureJWT.Claims{
				StandardClaims: jwt.StandardClaims{
					Subject:  userID,
					IssuedAt: expected.Unix(),
					ExpiresAt: func() int64 {
						n := expected
						return n.Add(gocaptureJWT.DefaultJWTExpirationDelta).Unix()
					}(),
				},
			},
		},
		{
			"claims with 30 min delta",
			1800, // 30min
			gocaptureJWT.Claims{
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
			var c gocaptureJWT.Claims
			// mock clock
			gocaptureJWT.SetClockInstance(&c, mockClock)
			c.Subject = userID
			c.IssueIt()
			c.SetExpirationDate(0)

			assert.Equal(t, userID, c.Subject)
		})
	}
}

func TestNewClaims(t *testing.T) {
	t.Parallel()

	c := gocaptureJWT.NewClaims("123", 0)
	err := c.Valid()
	assert.Nil(t, err)
}
