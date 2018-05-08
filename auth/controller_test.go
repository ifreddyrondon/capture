package auth_test

import (
	"net/http"
	"testing"

	"github.com/ifreddyrondon/gocapture/app"
	"github.com/ifreddyrondon/gocapture/auth"

	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/render/json"
)

func setupController(t *testing.T) (*bastion.Bastion, func()) {
	service, serviceTeardown := setupService(t)
	teardown := func() { serviceTeardown() }

	strategies := auth.Strategies{
		Service: service,
		Render:  json.NewRender,
		CtxKey:  app.ContextKey("user"),
	}

	controller := auth.Controller{
		Render:     json.NewRender,
		Strategies: &strategies,
	}

	app := bastion.New(bastion.Options{})
	app.APIRouter.Mount("/auth/", controller.Router())

	return app, teardown
}

func TestTokenAuthFailure(t *testing.T) {
	app, teardown := setupController(t)
	defer teardown()

	e := bastion.Tester(t, app)
	tt := []struct {
		name     string
		payload  map[string]interface{}
		response map[string]interface{}
	}{
		{
			name:    "invalid credentials",
			payload: map[string]interface{}{"email": testUserEmail, "password": "123"},
			response: map[string]interface{}{
				"status":  401.0,
				"error":   "Unauthorized",
				"message": "invalid email or password",
			},
		},
		{
			name:    "missing email",
			payload: map[string]interface{}{"email": "bla@localhost.com", "password": "123"},
			response: map[string]interface{}{
				"status":  401.0,
				"error":   "Unauthorized",
				"message": "invalid email or password",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			e.POST("/auth/token-auth").
				WithJSON(tc.payload).
				Expect().
				Status(http.StatusUnauthorized).
				JSON().Object().Equal(tc.response)
		})
	}
}

func TestTokenAuthFailureBadRequest(t *testing.T) {
	app, teardown := setupController(t)
	defer teardown()

	e := bastion.Tester(t, app)
	tc := struct {
		payload  map[string]interface{}
		response map[string]interface{}
	}{
		payload: map[string]interface{}{},
		response: map[string]interface{}{
			"status":  400.0,
			"error":   "Bad Request",
			"message": "email must not be blank\npassword must not be blank",
		},
	}

	e.POST("/auth/token-auth").
		WithJSON(tc.payload).
		Expect().
		Status(http.StatusBadRequest).
		JSON().Object().
		ContainsKey("status").ValueEqual("status", tc.response["status"]).
		ContainsKey("error").ValueEqual("error", tc.response["error"]).
		ContainsKey("message")
}
