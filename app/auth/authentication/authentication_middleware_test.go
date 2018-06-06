package authentication_test

import (
	"errors"
	"net/http"
	"sync"
	"testing"

	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/render/json"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/require"

	"github.com/ifreddyrondon/capture/app/auth/authentication"
	"github.com/ifreddyrondon/capture/app/auth/authentication/strategy/basic"
	"github.com/ifreddyrondon/capture/app/user"
	"github.com/ifreddyrondon/capture/database"
	"github.com/jinzhu/gorm"
)

var (
	once sync.Once
	db   *gorm.DB

	handler = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("OK"))
	})
)

const (
	testUserEmail    = "test@example.com"
	testUserPassword = "b4KeHAYy3u9v=ZQX"
)

func getDB() *gorm.DB {
	once.Do(func() {
		ds := database.Open("postgres://localhost/captures_app_test?sslmode=disable")
		db = ds.DB
	})
	return db
}

func setup(t *testing.T) (*bastion.Bastion, func()) {
	userStore := user.NewPGStore(getDB().Table("basic_auth-users"))
	userStore.Migrate()
	teardown := func() { userStore.Drop() }
	userService := user.NewService(userStore)

	// save a user to test
	u := user.User{Email: testUserEmail}
	err := u.SetPassword(testUserPassword)
	require.Nil(t, err)
	userService.Save(&u)

	strategy := basic.New(userService)
	middleware := authentication.NewAuthentication(strategy, json.NewRender)

	app := bastion.New(bastion.Options{})
	app.APIRouter.Route("/", func(r chi.Router) {
		r.Use(middleware.Authenticate)
		r.Post("/", handler)
	})

	return app, teardown
}

func TestAuthenticateSuccess(t *testing.T) {
	app, teardown := setup(t)
	defer teardown()

	payload := map[string]interface{}{"email": testUserEmail, "password": testUserPassword}
	e := bastion.Tester(t, app)
	e.POST("/").WithJSON(payload).
		Expect().
		Status(http.StatusOK).
		Text().Equal("OK")
}

func TestTokenAuthFailure(t *testing.T) {
	app, teardown := setup(t)
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
			payload: map[string]interface{}{"email": "bla@example.com", "password": "123"},
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

func TestTokenAuthFailureBadRequestCredentials(t *testing.T) {
	app, teardown := setup(t)
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

func setupWithMockStrategy(mock authentication.Strategy) *bastion.Bastion {
	middleware := authentication.NewAuthentication(mock, json.NewRender)

	app := bastion.New(bastion.Options{})
	app.APIRouter.Route("/", func(r chi.Router) {
		r.Use(middleware.Authenticate)
		r.Post("/", handler)
	})

	return app
}

type mockStrategyFailValidate struct{}

func (m *mockStrategyFailValidate) Validate(payload []byte) (*user.User, error) {
	return nil, errors.New("test")
}
func (m *mockStrategyFailValidate) IsErrCredentials(err error) bool {
	return false
}
func (m *mockStrategyFailValidate) IsErrDecoding(err error) bool {
	return false
}

func TestTokenAuthFailureBadRequestJSON(t *testing.T) {
	app := setupWithMockStrategy(&mockStrategyFailValidate{})

	e := bastion.Tester(t, app)
	tc := struct {
		payload  []byte
		response map[string]interface{}
	}{
		payload: []byte("{"),
		response: map[string]interface{}{
			"status":  400.0,
			"error":   "Bad Request",
			"message": "unexpected EOF",
		},
	}

	e.POST("/auth/token-auth").
		WithBytes(tc.payload).
		Expect().
		Status(http.StatusBadRequest).
		JSON().Object().Equal(tc.response)
}

func TestTokenAuthFailureInternalServerError(t *testing.T) {
	app := setupWithMockStrategy(&mockStrategyFailValidate{})

	e := bastion.Tester(t, app)
	tc := struct {
		payload  map[string]interface{}
		response map[string]interface{}
	}{
		payload: map[string]interface{}{"email": testUserEmail, "password": testUserPassword},
		response: map[string]interface{}{
			"status":  500.0,
			"error":   "Internal Server Error",
			"message": "test",
		},
	}

	e.POST("/auth/token-auth").
		WithJSON(tc.payload).
		Expect().
		Status(http.StatusInternalServerError).
		JSON().Object().Equal(tc.response)
}
