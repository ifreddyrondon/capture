package authorizing_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/ifreddyrondon/capture/pkg/authorizing"
	"github.com/ifreddyrondon/capture/pkg/domain"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"gopkg.in/src-d/go-kallax.v1"
)

type mockTokenService struct {
	subjectID string
	err       error
}

func (m *mockTokenService) IsRequestAuthorized(*http.Request) (string, error) {
	return m.subjectID, m.err
}

type mockStore struct {
	usr *domain.User
	err error
}

func (m *mockStore) GetUserByID(kallax.ULID) (*domain.User, error) { return m.usr, m.err }

func TestServiceAuthorizeRequest(t *testing.T) {
	t.Parallel()

	userIDTxt := "0162eb39-a65e-04a1-7ad9-d663bb49a396"
	userID, err := kallax.NewULIDFromText(userIDTxt)
	assert.Nil(t, err)
	u := &domain.User{ID: userID}

	s := authorizing.NewService(&mockTokenService{subjectID: userIDTxt}, &mockStore{usr: u})
	req, _ := http.NewRequest("GET", "/", nil)

	req.Header.Set("Authorization", "Bearer test")
	result, err := s.AuthorizeRequest(req)
	assert.Nil(t, err)
	assert.Equal(t, u, result)
}

func TestServiceAuthorizeRequestGetTokenFails(t *testing.T) {
	t.Parallel()

	s := authorizing.NewService(&mockTokenService{err: errors.New("test")}, &mockStore{})
	req, _ := http.NewRequest("GET", "/", nil)

	req.Header.Set("Authorization", "Bearer test")
	_, err := s.AuthorizeRequest(req)
	assert.Error(t, err)
}

type invalidErr interface{ IsInvalid() bool }

func TestServiceAuthorizeRequestInvalidSubjectID(t *testing.T) {
	t.Parallel()

	s := authorizing.NewService(&mockTokenService{subjectID: "a"}, &mockStore{})
	req, _ := http.NewRequest("GET", "/", nil)

	req.Header.Set("Authorization", "Bearer test")
	_, err := s.AuthorizeRequest(req)
	assert.EqualError(t, err, "a is not a valid ULID")
	authErr, ok := errors.Cause(err).(invalidErr)
	assert.True(t, ok)
	assert.True(t, authErr.IsInvalid())
}

type invalidCredentialErr string

func (i invalidCredentialErr) Error() string   { return fmt.Sprintf(string(i)) }
func (i invalidCredentialErr) IsInvalid() bool { return true }

type authorizationErr interface{ IsNotAuthorized() bool }

func TestServiceAuthorizeRequestInvalidCredentials(t *testing.T) {
	t.Parallel()

	ts := &mockTokenService{subjectID: "0162eb39-a65e-04a1-7ad9-d663bb49a396"}
	s := authorizing.NewService(ts, &mockStore{err: invalidCredentialErr("test")})
	req, _ := http.NewRequest("GET", "/", nil)

	req.Header.Set("Authorization", "Bearer test")
	_, err := s.AuthorizeRequest(req)
	assert.EqualError(t, err, "test")
	authErr, ok := errors.Cause(err).(authorizationErr)
	assert.True(t, ok)
	assert.True(t, authErr.IsNotAuthorized())
}
