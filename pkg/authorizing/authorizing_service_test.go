package authorizing_test

import (
	"net/http"
	"testing"

	"github.com/ifreddyrondon/capture/pkg"
	"github.com/ifreddyrondon/capture/pkg/authorizing"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"gopkg.in/src-d/go-kallax.v1"
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

func (m *mockStore) GetUserByID(kallax.ULID) (*pkg.User, error) { return m.usr, m.err }

func TestServiceAuthorizeRequest(t *testing.T) {
	t.Parallel()

	uidText := "0162eb39-a65e-04a1-7ad9-d663bb49a396"
	uid, err := kallax.NewULIDFromText(uidText)
	assert.Nil(t, err)
	u := &pkg.User{ID: uid}

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

func TestServiceAuthorizeRequestInvalidSubjectID(t *testing.T) {
	t.Parallel()

	s := authorizing.NewService(&mockTokenService{subjectID: "a"}, &mockStore{})
	req, _ := http.NewRequest("GET", "/", nil)

	req.Header.Set("Authorization", "Bearer test")
	_, err := s.AuthorizeRequest(req)
	assert.EqualError(t, err, "uuid: UUID string too short: a")
	authErr, ok := errors.Cause(err).(authorizationErr)
	assert.True(t, ok)
	assert.True(t, authErr.IsNotAuthorized())
}
