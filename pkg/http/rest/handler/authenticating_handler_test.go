package handler_test

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/ifreddyrondon/capture/pkg/http/rest/handler"

	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/capture/pkg/authenticating"
	"github.com/ifreddyrondon/capture/pkg/domain"
	"gopkg.in/src-d/go-kallax.v1"
)

type mockAuthenticatingService struct {
	usr      *domain.User
	token    string
	err      error
	tokenErr error
}

func (s *mockAuthenticatingService) AuthenticateUser(credential authenticating.BasicCredential) (*domain.User, error) {
	return s.usr, s.err
}
func (s *mockAuthenticatingService) GetUserToken(kallax.ULID) (string, error) {
	return s.token, s.tokenErr
}

type invalidCredentialErr string

func (i invalidCredentialErr) Error() string            { return fmt.Sprintf(string(i)) }
func (i invalidCredentialErr) InvalidCredentials() bool { return true }

func TestAuthenticateSuccess(t *testing.T) {
	t.Parallel()

	s := &mockAuthenticatingService{
		usr:   &domain.User{ID: kallax.NewULID()},
		token: "token*test",
	}
	app := bastion.New()
	app.APIRouter.Post("/", handler.Authenticating(s))

	response := map[string]interface{}{"token": "token*test"}
	e := bastion.Tester(t, app)
	payload := map[string]interface{}{"email": "bla@example.com", "password": "123"}
	e.POST("/").WithJSON(payload).
		Expect().
		Status(http.StatusOK).
		JSON().Object().Equal(response)
}

func TestAuthenticateFailBadRequest(t *testing.T) {
	t.Parallel()

	s := &mockAuthenticatingService{usr: &domain.User{}}
	app := bastion.New()
	app.APIRouter.Post("/", handler.Authenticating(s))

	response := map[string]interface{}{
		"status":  400.0,
		"error":   "Bad Request",
		"message": "invalid email",
	}
	e := bastion.Tester(t, app)
	payload := map[string]interface{}{"email": "bla@example", "password": "123"}
	e.POST("/").WithJSON(payload).
		Expect().
		Status(http.StatusBadRequest).
		JSON().Object().Equal(response)
}

func TestAuthenticateFailUnauthorized(t *testing.T) {
	t.Parallel()

	s := &mockAuthenticatingService{err: invalidCredentialErr("invalid email or password")}
	app := bastion.New()
	app.APIRouter.Post("/", handler.Authenticating(s))

	response := map[string]interface{}{
		"status":  401.0,
		"error":   "Unauthorized",
		"message": "invalid email or password",
	}
	e := bastion.Tester(t, app)
	payload := map[string]interface{}{"email": "bla@example.com", "password": "123"}
	e.POST("/").WithJSON(payload).
		Expect().
		Status(http.StatusUnauthorized).
		JSON().Object().Equal(response)
}

func TestAuthenticateFailInternalServerErrorWhenAuthenticateUser(t *testing.T) {
	t.Parallel()

	s := &mockAuthenticatingService{err: errors.New("test")}
	app := bastion.New()
	app.APIRouter.Post("/", handler.Authenticating(s))

	response := map[string]interface{}{
		"status":  500.0,
		"error":   "Internal Server Error",
		"message": "looks like something went wrong",
	}
	e := bastion.Tester(t, app)
	payload := map[string]interface{}{"email": "bla@example.com", "password": "123"}
	e.POST("/").WithJSON(payload).
		Expect().
		Status(http.StatusInternalServerError).
		JSON().Object().Equal(response)
}

func TestAuthenticateFailInternalServerErrorWhenGetUserToken(t *testing.T) {
	t.Parallel()

	s := &mockAuthenticatingService{usr: &domain.User{}, tokenErr: errors.New("test")}
	app := bastion.New()
	app.APIRouter.Post("/", handler.Authenticating(s))

	response := map[string]interface{}{
		"status":  500.0,
		"error":   "Internal Server Error",
		"message": "looks like something went wrong",
	}
	e := bastion.Tester(t, app)
	payload := map[string]interface{}{"email": "bla@example.com", "password": "123"}
	e.POST("/").WithJSON(payload).
		Expect().
		Status(http.StatusInternalServerError).
		JSON().Object().Equal(response)
}
