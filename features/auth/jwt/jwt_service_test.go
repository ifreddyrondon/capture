package jwt_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ifreddyrondon/capture/features/auth/jwt"
)

// setup a valid service
func getValidService() *jwt.Service {
	return jwt.NewService([]byte("secret"), jwt.DefaultJWTExpirationDelta)
}

func TestServiceGenerateToken(t *testing.T) {
	s := getValidService()

	tt := []struct {
		name   string
		userID string
	}{
		{"valid", "123"},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			token, err := s.GenerateToken(tc.userID)
			assert.Nil(t, err)
			assert.NotEmpty(t, token)
		})
	}
}

func TestAuthorizationFromHeader(t *testing.T) {
	t.Parallel()

	subjectID := "test_123"
	s := getValidService()
	token, err := s.GenerateToken(subjectID)
	assert.Nil(t, err)

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
	subj, err := s.IsAuthorizedREQ(req)
	assert.Nil(t, err)
	assert.Equal(t, subjectID, subj)
}

func TestAuthorizationPostForm(t *testing.T) {
	t.Parallel()

	subjectID := "test_123"
	s := getValidService()
	token, err := s.GenerateToken(subjectID)
	assert.Nil(t, err)

	data := url.Values{}
	data.Add("access_token", token)

	req := &http.Request{
		Method: "POST",
		Header: http.Header{"Content-Type": {`application/x-www-form-urlencoded`}},
		Body:   ioutil.NopCloser(strings.NewReader(data.Encode())),
	}

	subj, err := s.IsAuthorizedREQ(req)
	assert.Nil(t, err)
	assert.Equal(t, subjectID, subj)
}

func TestAuthorizationFailInvalidToken(t *testing.T) {
	t.Parallel()

	s := getValidService()
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c")
	_, err := s.IsAuthorizedREQ(req)
	assert.EqualError(t, err, "you don’t have permission to access this resource")
	assert.True(t, s.IsNotAuthorizedErr(err))
}

func TestAuthorizationFailInvalidSignedMethod(t *testing.T) {
	t.Parallel()

	s := getValidService()
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.TCYt5XsITJX1CxPCT8yAV-TVkIEq_PbChOMqsLfRoPsnsgw5WEuts01mq-pQy7UJiN5mgRxD-WUcX16dUEMGlv50aqzpqh4Qktb3rk-BuQy72IFLOqV0G_zS245-kronKb78cPN25DGlcTwLtjPAYuNzVBAh4vGHSrQyHUdBBPM")
	_, err := s.IsAuthorizedREQ(req)
	assert.EqualError(t, err, "you don’t have permission to access this resource")
}
