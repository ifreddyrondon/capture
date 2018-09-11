package user_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/go-chi/chi"
	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/capture/features/auth/authorization"
	"github.com/ifreddyrondon/capture/features/user"
	"github.com/stretchr/testify/require"
	"gopkg.in/src-d/go-kallax.v1"
)

const (
	testUserEmail    = "test@example.com"
	testUserPassword = "b4KeHAYy3u9v=ZQX"
)

func setCtxMiddleware(subjectID string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := authorization.WithSubjectID(r.Context(), subjectID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

var handler = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(testUserEmail))
})

func setupMiddleware(subjectID string, service user.GetterService) *bastion.Bastion {
	app := bastion.New()
	app.APIRouter.Route("/", func(r chi.Router) {
		r.Use(setCtxMiddleware(subjectID))
		r.Use(user.LoggedUser(service))
		r.Get("/", handler)
	})

	return app
}

func TestLoggedUserMiddlewareBadRequestInvalidUUID(t *testing.T) {
	t.Parallel()

	app := setupMiddleware("1", &user.MockService{})
	e := bastion.Tester(t, app)
	tc := struct {
		response map[string]interface{}
	}{
		response: map[string]interface{}{
			"status":  400.0,
			"error":   "Bad Request",
			"message": "uuid: UUID string too short: 1",
		},
	}

	e.GET("/").
		Expect().
		Status(http.StatusBadRequest).
		JSON().Object().Equal(tc.response)
}

func TestLoggedUserMiddlewareInternalServerError(t *testing.T) {
	t.Parallel()

	app := setupMiddleware(kallax.NewULID().String(), &user.MockService{Err: errors.New("test")})

	e := bastion.Tester(t, app)
	tc := struct {
		response map[string]interface{}
	}{
		response: map[string]interface{}{
			"status":  500.0,
			"error":   "Internal Server Error",
			"message": "looks like something went wrong",
		},
	}

	e.GET("/").
		Expect().
		Status(http.StatusInternalServerError).
		JSON().Object().Equal(tc.response)
}

func setupServiceMiddleware(t *testing.T) (string, user.GetterService, func()) {
	service, teardown := setupService(t)
	u := user.User{Email: testUserEmail}
	err := u.SetPassword(testUserPassword)
	require.Nil(t, err)
	service.Save(&u)

	return u.ID.String(), service, teardown
}

func TestLoggedUserMiddlewareOK(t *testing.T) {
	t.Parallel()

	subjectID, service, teardown := setupServiceMiddleware(t)
	defer teardown()
	app := setupMiddleware(subjectID, service)

	e := bastion.Tester(t, app)
	e.GET("/").
		Expect().
		Status(http.StatusOK).
		Text().Equal(testUserEmail)
}

func TestLoggedUserMiddlewareUserNotFound(t *testing.T) {
	t.Parallel()

	_, service, teardown := setupServiceMiddleware(t)
	defer teardown()
	// set another subjectID and force not found user
	app := setupMiddleware(kallax.NewULID().String(), service)

	e := bastion.Tester(t, app)
	tc := struct {
		response map[string]interface{}
	}{
		response: map[string]interface{}{
			"status":  404.0,
			"error":   "Not Found",
			"message": "user not found",
		},
	}

	e.GET("/").
		Expect().
		Status(http.StatusNotFound).
		JSON().Object().Equal(tc.response)
}
