package authenticating_test

import (
	"fmt"
	"testing"

	"github.com/ifreddyrondon/capture/pkg/authenticating"
	"github.com/ifreddyrondon/capture/pkg/domain"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

type authenticatingErr interface{ InvalidCredentials() bool }

type userNotFoundMock string

func (u userNotFoundMock) Error() string  { return fmt.Sprint("test") }
func (u userNotFoundMock) NotFound() bool { return true }

type mockStore struct {
	usr *domain.User
	err error
}

func (m *mockStore) GetUserByEmail(string) (*domain.User, error) { return m.usr, m.err }

type mockTokenService struct {
	token string
	err   error
}

func (m *mockTokenService) GenerateToken(string) (string, error) { return m.token, m.err }

func TestAuthenticatingServiceGenerateToken(t *testing.T) {
	t.Parallel()

	s := authenticating.NewService(&mockTokenService{token: "a"}, &mockStore{})
	result, err := s.GetUserToken("a")
	assert.Nil(t, err)
	assert.Equal(t, "a", result)
}

func TestAuthenticatingServiceGenerateTokenErr(t *testing.T) {
	t.Parallel()

	s := authenticating.NewService(&mockTokenService{token: "", err: errors.New("test")}, &mockStore{})
	_, err := s.GetUserToken("a")
	assert.Error(t, err)
}

func TestAuthenticatingService(t *testing.T) {
	t.Parallel()

	credential := authenticating.BasicCredential{Password: "secret"}
	mockUser := &domain.User{Password: []byte("$2a$14$ajq8Q7fbtFRQvXpdCq7Jcuy.Rx1h/L4J60Otx.gyNLbAYctGMJ9tK")}
	s := authenticating.NewService(&mockTokenService{}, &mockStore{usr: mockUser})
	result, err := s.AuthenticateUser(credential)
	assert.Nil(t, err)
	assert.Equal(t, mockUser, result)
}

func TestAuthenticatingServiceFailWhenUserNotFound(t *testing.T) {
	t.Parallel()
	s := authenticating.NewService(&mockTokenService{}, &mockStore{err: userNotFoundMock("")})
	_, err := s.AuthenticateUser(authenticating.BasicCredential{})
	assert.EqualError(t, err, "test")
	authErr, ok := errors.Cause(err).(authenticatingErr)
	assert.True(t, ok)
	assert.True(t, authErr.InvalidCredentials())

}

func TestAuthenticatingServiceFailError(t *testing.T) {
	t.Parallel()
	s := authenticating.NewService(&mockTokenService{}, &mockStore{err: errors.New("test")})
	_, err := s.AuthenticateUser(authenticating.BasicCredential{})
	assert.EqualError(t, err, "test")
}

func TestAuthenticatingServiceFailInvalidPassword(t *testing.T) {
	t.Parallel()

	credential := authenticating.BasicCredential{Password: "secret2"}
	mockUser := &domain.User{Password: []byte("$2a$14$ajq8Q7fbtFRQvXpdCq7Jcuy.Rx1h/L4J60Otx.gyNLbAYctGMJ9tK")}
	s := authenticating.NewService(&mockTokenService{}, &mockStore{usr: mockUser})
	_, err := s.AuthenticateUser(credential)
	assert.EqualError(t, err, "invalid password")
	authErr, ok := errors.Cause(err).(authenticatingErr)
	assert.True(t, ok)
	assert.True(t, authErr.InvalidCredentials())
}
