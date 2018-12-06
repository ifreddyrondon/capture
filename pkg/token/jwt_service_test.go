package token_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/ifreddyrondon/capture/pkg/token"
)

// setup a valid service
func getValidService() *token.JwtService {
	return token.NewJWTService("secret", time.Minute)
}

type authorizationErr interface{ IsNotAuthorized() bool }

func TestServiceGenerateToken(t *testing.T) {
	t.Parallel()
	s := getValidService()

	tt := []struct {
		name   string
		userID string
	}{
		{"valid", "123"},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			hash, err := s.GenerateToken(tc.userID)
			assert.Nil(t, err)
			assert.NotEmpty(t, hash)
		})
	}
}

func TestAuthorizingFromHeader(t *testing.T) {
	t.Parallel()

	subjectID := "test_123"
	s := getValidService()
	tok, err := s.GenerateToken(subjectID)
	assert.Nil(t, err)

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tok))
	subj, err := s.IsRequestAuthorized(req)
	assert.Nil(t, err)
	assert.Equal(t, subjectID, subj)
}

func TestAuthorizingPostForm(t *testing.T) {
	t.Parallel()

	subjectID := "test_123"
	s := getValidService()
	tok, err := s.GenerateToken(subjectID)
	assert.Nil(t, err)

	data := url.Values{}
	data.Add("access_token", tok)

	req := &http.Request{
		Method: "POST",
		Header: http.Header{"Content-Type": {`application/x-www-form-urlencoded`}},
		Body:   ioutil.NopCloser(strings.NewReader(data.Encode())),
	}

	subj, err := s.IsRequestAuthorized(req)
	assert.Nil(t, err)
	assert.Equal(t, subjectID, subj)
}

func TestAuthorizingFailInvalidToken(t *testing.T) {
	t.Parallel()

	s := getValidService()
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c")
	_, err := s.IsRequestAuthorized(req)
	assert.EqualError(t, err, "signature is invalid")
	authErr, ok := errors.Cause(err).(authorizationErr)
	assert.True(t, ok)
	assert.True(t, authErr.IsNotAuthorized())
}

func TestAuthorizingFailInvalidSignedMethod(t *testing.T) {
	t.Parallel()

	s := getValidService()
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.TCYt5XsITJX1CxPCT8yAV-TVkIEq_PbChOMqsLfRoPsnsgw5WEuts01mq-pQy7UJiN5mgRxD-WUcX16dUEMGlv50aqzpqh4Qktb3rk-BuQy72IFLOqV0G_zS245-kronKb78cPN25DGlcTwLtjPAYuNzVBAh4vGHSrQyHUdBBPM")
	_, err := s.IsRequestAuthorized(req)
	assert.EqualError(t, err, "unexpected signing method")
	authErr, ok := errors.Cause(err).(authorizationErr)
	assert.True(t, ok)
	assert.True(t, authErr.IsNotAuthorized())
}
