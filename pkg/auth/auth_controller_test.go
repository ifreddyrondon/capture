package auth_test

import (
	"net/http"
	"testing"

	"github.com/ifreddyrondon/capture/pkg"
	"github.com/ifreddyrondon/capture/pkg/auth"
	"github.com/ifreddyrondon/capture/pkg/auth/jwt"
	"github.com/ifreddyrondon/capture/pkg/user"
	"gopkg.in/src-d/go-kallax.v1"

	"github.com/ifreddyrondon/bastion"
)

const (
	userEmail    = "test@example.com"
	userPassword = "b4KeHAYy3u9v=ZQX"
)

func notLoggedUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func loggedUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var mockUser = &pkg.User{Email: "test@example.com", ID: kallax.NewULID()}
		ctx := user.WithUser(r.Context(), mockUser)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func setup(middleware func(next http.Handler) http.Handler) *bastion.Bastion {
	jwtService := jwt.NewService([]byte("test"), jwt.DefaultJWTExpirationDelta)

	app := bastion.New()
	app.APIRouter.Mount("/auth/", auth.Routes(middleware, jwtService))

	return app
}

func TestBasicAuthenticationOK(t *testing.T) {
	app := setup(loggedUser)

	payload := map[string]interface{}{"email": userEmail, "password": userPassword}
	e := bastion.Tester(t, app)
	e.POST("/auth/token-auth").WithJSON(payload).
		Expect().
		Status(http.StatusOK).
		JSON().Object().
		ContainsKey("token")
}

func TestBasicAuthenticationMissUser(t *testing.T) {
	app := setup(notLoggedUser)

	payload := map[string]interface{}{"email": userEmail, "password": userPassword}
	response := map[string]interface{}{
		"status":  500.0,
		"error":   "Internal Server Error",
		"message": "looks like something went wrong",
	}
	e := bastion.Tester(t, app)
	e.POST("/auth/token-auth").WithJSON(payload).
		Expect().
		Status(http.StatusInternalServerError).
		JSON().Object().Equal(response)
}
