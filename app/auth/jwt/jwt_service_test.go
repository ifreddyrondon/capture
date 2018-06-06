package jwt_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"

	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/render/json"
	"github.com/ifreddyrondon/capture/app/auth/jwt"
)

func TestServiceGenerateToken(t *testing.T) {
	s := jwt.NewService([]byte("secret"), jwt.DefaultJWTExpirationDelta, json.NewRender)

	tt := []struct {
		name   string
		userID string
	}{
		{"valid", "123"},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			token, err := s.GenerateToken(tc.userID)
			assert.Nil(t, err)
			assert.NotEmpty(t, token)
		})
	}
}

// setup a valid service
func getValidService() *jwt.Service {
	return jwt.NewService([]byte("secret"), jwt.DefaultJWTExpirationDelta, json.NewRender)
}

func setupApp(authorizationMidleware func(http.Handler) http.Handler) *bastion.Bastion {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "ok")
	})
	app := bastion.New(bastion.Options{})
	app.APIRouter.Route("/", func(r chi.Router) {
		r.Use(authorizationMidleware)
		r.Get("/", handler)
		r.Post("/", handler)
	})
	return app
}

func TestAuthorizationFromHeader(t *testing.T) {
	t.Parallel()

	s := getValidService()
	app := setupApp(s.IsAuthorized)

	token, err := s.GenerateToken("123")
	assert.Nil(t, err)
	e := bastion.Tester(t, app)
	e.GET("/").WithHeader("Authorization", fmt.Sprintf("Bearer %v", token)).
		Expect().
		Status(http.StatusOK)
}

func TestAuthorizationPostForm(t *testing.T) {
	t.Parallel()

	s := getValidService()
	app := setupApp(s.IsAuthorized)

	token, err := s.GenerateToken("123")
	assert.Nil(t, err)
	e := bastion.Tester(t, app)
	e.POST("/").WithForm(map[string]string{"access_token": token}).
		Expect().
		Status(http.StatusOK)
}

func TestAuthorizationFailInvalidToken(t *testing.T) {
	t.Parallel()

	response := map[string]interface{}{
		"status":  403.0,
		"error":   "Forbidden",
		"message": "you don’t have permission to access this resource",
	}

	s := getValidService()
	app := setupApp(s.IsAuthorized)
	e := bastion.Tester(t, app)
	e.GET("/").WithHeader("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c").
		Expect().
		Status(http.StatusForbidden).
		JSON().Object().Equal(response)
}

func TestAuthorizationFailInvalidSignedMethod(t *testing.T) {
	t.Parallel()
	// the return err is the same but internaly validates signed method first
	response := map[string]interface{}{
		"status":  403.0,
		"error":   "Forbidden",
		"message": "you don’t have permission to access this resource",
	}

	s := getValidService()
	app := setupApp(s.IsAuthorized)
	e := bastion.Tester(t, app)
	e.GET("/").WithHeader("Authorization", "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.TCYt5XsITJX1CxPCT8yAV-TVkIEq_PbChOMqsLfRoPsnsgw5WEuts01mq-pQy7UJiN5mgRxD-WUcX16dUEMGlv50aqzpqh4Qktb3rk-BuQy72IFLOqV0G_zS245-kronKb78cPN25DGlcTwLtjPAYuNzVBAh4vGHSrQyHUdBBPM").
		Expect().
		Status(http.StatusForbidden).
		JSON().Object().Equal(response)
}
