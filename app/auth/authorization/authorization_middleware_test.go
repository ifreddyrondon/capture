package authorization_test

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/go-chi/chi"
	"github.com/ifreddyrondon/bastion"
	"github.com/stretchr/testify/assert"

	"github.com/ifreddyrondon/capture/app/auth/authorization"
	"github.com/ifreddyrondon/capture/app/auth/jwt"
)

// setup a valid service
var jwtService = jwt.NewService([]byte("secret"), jwt.DefaultJWTExpirationDelta)

func setupApp(authorizationStrategy authorization.Strategy) *bastion.Bastion {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "ok")
	})
	middle := authorization.NewAuthorization(authorizationStrategy)

	app := bastion.New()
	app.APIRouter.Route("/", func(r chi.Router) {
		r.Use(middle.IsAuthorizedREQ)
		r.Get("/", handler)
		r.Post("/", handler)
	})
	return app
}

func TestAuthorizationFromHeader(t *testing.T) {
	t.Parallel()

	app := setupApp(jwtService)

	token, err := jwtService.GenerateToken("123")
	assert.Nil(t, err)
	e := bastion.Tester(t, app)
	e.GET("/").WithHeader("Authorization", fmt.Sprintf("Bearer %v", token)).
		Expect().
		Status(http.StatusOK)
}

func TestAuthorizationPostForm(t *testing.T) {
	t.Parallel()

	app := setupApp(jwtService)

	token, err := jwtService.GenerateToken("123")
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

	app := setupApp(jwtService)
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

	app := setupApp(jwtService)
	e := bastion.Tester(t, app)
	e.GET("/").WithHeader("Authorization", "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.TCYt5XsITJX1CxPCT8yAV-TVkIEq_PbChOMqsLfRoPsnsgw5WEuts01mq-pQy7UJiN5mgRxD-WUcX16dUEMGlv50aqzpqh4Qktb3rk-BuQy72IFLOqV0G_zS245-kronKb78cPN25DGlcTwLtjPAYuNzVBAh4vGHSrQyHUdBBPM").
		Expect().
		Status(http.StatusForbidden).
		JSON().Object().Equal(response)
}

type mockStrategyFailIsAuthorizedREQ struct{}

func (m *mockStrategyFailIsAuthorizedREQ) IsAuthorizedREQ(r *http.Request) (string, error) {
	return "", errors.New("test")
}
func (m *mockStrategyFailIsAuthorizedREQ) IsNotAuthorizedErr(err error) bool {
	return false
}

func TestAuthorizationFailInternalServerError(t *testing.T) {
	t.Parallel()

	app := setupApp(&mockStrategyFailIsAuthorizedREQ{})
	response := map[string]interface{}{
		"status":  500.0,
		"error":   "Internal Server Error",
		"message": "looks like something went wrong",
	}

	token, err := jwtService.GenerateToken("123")
	assert.Nil(t, err)
	e := bastion.Tester(t, app)
	e.GET("/").WithHeader("Authorization", fmt.Sprintf("Bearer %v", token)).
		Expect().
		Status(http.StatusInternalServerError).
		JSON().Object().Equal(response)
}
