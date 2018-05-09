package jwt_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ifreddyrondon/gocapture/jwt"
)

func TestServiceGenerateToken(t *testing.T) {
	s := jwt.NewService([]byte("secret"), jwt.DefaultJWTExpirationDelta)

	tt := []struct {
		name   string
		userID string
	}{
		{
			"valid",
			"123",
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			token, err := s.GenerateToken(tc.userID)
			assert.Nil(t, err)
			assert.NotEmpty(t, token)
		})
	}
}
