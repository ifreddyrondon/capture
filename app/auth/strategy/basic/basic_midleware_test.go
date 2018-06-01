package basic_test

import (
	"net/http"
	"sync"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/require"

	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/render/json"

	"github.com/ifreddyrondon/capture/app/auth/strategy/basic"
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

	service := basic.NewStrategy(json.NewRender, userService)
	app := bastion.New(bastion.Options{})
	app.APIRouter.Route("/", func(r chi.Router) {
		r.Use(service.Authenticate)
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

func TestTokenAuthFailureBadRequest(t *testing.T) {
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
