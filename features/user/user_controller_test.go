package user_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/capture/features/user"
)

func setupController(service user.Service) *bastion.Bastion {
	app := bastion.New()
	app.APIRouter.Mount("/users/", user.Routes(service))

	return app
}

func TestCreateValidUser(t *testing.T) {
	service, teardown := setupService(t)
	defer teardown()
	app := setupController(service)

	e := bastion.Tester(t, app)
	tt := []struct {
		name     string
		payload  map[string]interface{}
		response map[string]interface{}
	}{
		{
			name:     "create user with only email",
			payload:  map[string]interface{}{"email": "test@example.com"},
			response: map[string]interface{}{"email": "test@example.com"},
		},
		{
			name: "create user",
			payload: map[string]interface{}{
				"email":    "test2@example.com",
				"password": "b4KeHAYy3u9v=ZQX",
			},
			response: map[string]interface{}{"email": "test2@example.com"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			e.POST("/users/").
				WithJSON(tc.payload).
				Expect().
				Status(http.StatusCreated).
				JSON().Object().
				ContainsKey("email").ValueEqual("email", tc.response["email"]).
				ContainsKey("id").NotEmpty().
				ContainsKey("createdAt").NotEmpty().
				ContainsKey("updatedAt").NotEmpty().
				NotContainsKey("password")
		})
	}
}

func TestCreateINValidUser(t *testing.T) {
	service, teardown := setupService(t)
	defer teardown()
	app := setupController(service)

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
			name:    "missing email",
			payload: map[string]interface{}{"email": ""},
			response: map[string]interface{}{
				"status":  400.0,
				"error":   "Bad Request",
				"message": "email must not be blank",
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
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			e.POST("/users/").
				WithJSON(tc.payload).
				Expect().
				Status(http.StatusBadRequest).
				JSON().Object().Equal(tc.response)
		})
	}
}

func TestConflictEmail(t *testing.T) {
	service, teardown := setupService(t)
	defer teardown()
	app := setupController(service)

	payload := map[string]interface{}{"email": "test@example.com"}
	response := map[string]interface{}{
		"status":  409.0,
		"error":   "Conflict",
		"message": "email 'test@example.com' already exists",
	}

	e := bastion.Tester(t, app)
	e.POST("/users/").WithJSON(payload).Expect().Status(http.StatusCreated)

	e.POST("/users/").WithJSON(payload).
		Expect().
		Status(http.StatusConflict).
		JSON().Object().Equal(response)
}

func TestCreateFailSave(t *testing.T) {
	t.Parallel()

	service := &user.MockService{Err: errors.New("test")}
	app := setupController(service)

	payload := map[string]interface{}{"email": "test@example.com"}
	response := map[string]interface{}{
		"status":  500.0,
		"error":   "Internal Server Error",
		"message": "looks like something went wrong",
	}

	e := bastion.Tester(t, app)
	e.POST("/users/").WithJSON(payload).
		Expect().
		Status(http.StatusInternalServerError).
		JSON().Object().Equal(response)
}
