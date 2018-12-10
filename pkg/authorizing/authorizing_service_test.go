package authorizing_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/ifreddyrondon/capture/pkg"
	"github.com/ifreddyrondon/capture/pkg/authorizing"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

type authorizationErr interface{ IsNotAuthorized() bool }

type mockTokenService struct {
	subjectID string
	err       error
}

func (m *mockTokenService) IsRequestAuthorized(*http.Request) (string, error) {
	return m.subjectID, m.err
}

type mockStore struct {
	usr *pkg.User
	err error
}

func (m *mockStore) GetUserByID(string) (*pkg.User, error) { return m.usr, m.err }

func TestServiceAuthorizeRequest(t *testing.T) {
	t.Parallel()

	uidText := "0162eb39-a65e-04a1-7ad9-d663bb49a396"
	u := &pkg.User{ID: uidText}

	s := authorizing.NewService(&mockTokenService{subjectID: uidText}, &mockStore{usr: u})
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

type invalidErr string

func (i invalidErr) Error() string   { return fmt.Sprintf(string(i)) }
func (i invalidErr) IsInvalid() bool { return true }

func TestServiceAuthorizeRequestInvalidSubjectID(t *testing.T) {
	t.Parallel()

	s := authorizing.NewService(&mockTokenService{subjectID: "a"}, &mockStore{err: invalidErr("test")})
	req, _ := http.NewRequest("GET", "/", nil)

	req.Header.Set("Authorization", "Bearer test")
	_, err := s.AuthorizeRequest(req)
	assert.EqualError(t, err, "test")
	authErr, ok := errors.Cause(err).(authorizationErr)
	assert.True(t, ok)
	assert.True(t, authErr.IsNotAuthorized())
}
