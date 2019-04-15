package handler_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/ifreddyrondon/capture/pkg/http/rest/handler"

	"github.com/ifreddyrondon/bastion"

	"github.com/ifreddyrondon/capture/pkg/signup"
)

type mockSignUpService struct {
	usr *signup.User
	err error
}

func (s *mockSignUpService) EnrollUser(payl signup.Payload) (*signup.User, error) { return s.usr, s.err }

func TestSignUpSuccess(t *testing.T) {
	t.Parallel()
	id := "0162eb39-a65e-04a1-7ad9-d663bb49a396"

	tt := []struct {
		name    string
		payload map[string]interface{}
		usrMock *signup.User
	}{
		{
			name:    "create user with only email",
			payload: map[string]interface{}{"email": "bla@example.com", "password": "1234"},
			usrMock: &signup.User{ID: id, Email: "bla@example.com"},
		},
		{
			name:    "create user",
			payload: map[string]interface{}{"email": "bla@example.com"},
			usrMock: &signup.User{ID: id, Email: "bla@example.com"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			s := &mockSignUpService{usr: tc.usrMock}
			app := bastion.New()
			app.Post("/", handler.SignUp(s))
			e := bastion.Tester(t, app)
			e.POST("/").WithJSON(tc.payload).
				Expect().
				Status(http.StatusCreated).
				JSON().Object().
				ContainsKey("email").ValueEqual("email", tc.payload["email"]).
				ContainsKey("id").ValueEqual("id", "0162eb39-a65e-04a1-7ad9-d663bb49a396").
				ContainsKey("createdAt").NotEmpty().
				ContainsKey("updatedAt").NotEmpty().
				NotContainsKey("password")
		})
	}
}

func TestSignUpBadRequest(t *testing.T) {
	t.Parallel()

	s := &mockSignUpService{}
	app := bastion.New()
	app.Post("/", handler.SignUp(s))

	e := bastion.Tester(t, app)
	tt := []struct {
		name     string
		payload  map[string]interface{}
		response map[string]interface{}
	}{
		{
			name:    "invalid email",
			payload: map[string]interface{}{"email": "test@asd"},
			response: map[string]interface{}{
				"status":  400.0,
				"error":   "Bad Request",
				"message": "invalid email",
			},
		},
		{
			name:    "missing email - empty",
			payload: map[string]interface{}{"email": ""},
			response: map[string]interface{}{
				"status":  400.0,
				"error":   "Bad Request",
				"message": "invalid email",
			},
		},
		{
			name:    "missing email",
			payload: map[string]interface{}{},
			response: map[string]interface{}{
				"status":  400.0,
				"error":   "Bad Request",
				"message": "email must not be blank",
			},
		},
		{
			name:    "invalid password",
			payload: map[string]interface{}{"email": "bla@example.com", "password": "1"},
			response: map[string]interface{}{
				"status":  400.0,
				"error":   "Bad Request",
				"message": "password must have at least four characters",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			e.POST("/").
				WithJSON(tc.payload).
				Expect().
				Status(http.StatusBadRequest).
				JSON().Object().Equal(tc.response)
		})
	}
}

type conflictErr string

func (e conflictErr) Error() string  { return string(e) }
func (e conflictErr) Conflict() bool { return true }

func TestSignUpConflictEmail(t *testing.T) {
	s := &mockSignUpService{err: conflictErr("test")}
	app := bastion.New()
	app.Post("/", handler.SignUp(s))

	payload := map[string]interface{}{"email": "test@example.com"}
	response := map[string]interface{}{
		"status":  409.0,
		"error":   "Conflict",
		"message": "email 'test@example.com' already exists",
	}

	e := bastion.Tester(t, app)
	e.POST("/").WithJSON(payload).
		Expect().
		Status(http.StatusConflict).
		JSON().Object().Equal(response)
}

func TestSignUpFailSave(t *testing.T) {
	t.Parallel()

	s := &mockSignUpService{err: errors.New("test")}
	app := bastion.New()
	app.Post("/", handler.SignUp(s))

	payload := map[string]interface{}{"email": "test@example.com"}
	response := map[string]interface{}{
		"status":  500.0,
		"error":   "Internal Server Error",
		"message": "looks like something went wrong",
	}

	e := bastion.Tester(t, app)
	e.POST("/").WithJSON(payload).
		Expect().
		Status(http.StatusInternalServerError).
		JSON().Object().Equal(response)
}
