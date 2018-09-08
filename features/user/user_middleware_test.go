package user_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/ifreddyrondon/capture/features/auth/authorization"
	"github.com/stretchr/testify/require"

	"github.com/go-chi/chi"
	"github.com/ifreddyrondon/bastion"

	"github.com/ifreddyrondon/capture/features/user"
	"gopkg.in/src-d/go-kallax.v1"
)

const (
	testUserEmail    = "test@example.com"
	testUserPassword = "b4KeHAYy3u9v=ZQX"
)

func setCtxMidd(subjectID string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := authorization.NewContextManager().WithSubjectID(r.Context(), subjectID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

var handler = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(testUserEmail))
})

type mockUserGetterServiceFail struct{}

func (m *mockUserGetterServiceFail) GetByEmail(email string) (*user.User, error) {
	return nil, errors.New("test")
}

func (m *mockUserGetterServiceFail) GetByID(id kallax.ULID) (*user.User, error) {
	return nil, errors.New("test")
}

func getNewMiddleware(service user.GetterService) *user.Middleware {
	return user.NewMiddleware(service)
}

func setupWithMockStrategy(subjectID string) *bastion.Bastion {
	middleware := getNewMiddleware(&mockUserGetterServiceFail{})

	app := bastion.New()
	app.APIRouter.Route("/", func(r chi.Router) {
		r.Use(setCtxMidd(subjectID))
		r.Use(middleware.LoggedUser)
		r.Get("/", handler)
	})

	return app
}

func TestLoggedUserMiddlewareInternalServerError(t *testing.T) {
	t.Parallel()

	app := setupWithMockStrategy("1")

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

func TestLoggedUserMiddlewareBadRequestInvalidUUID(t *testing.T) {
	t.Parallel()

	app := setupWithMockStrategy(kallax.NewULID().String())

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

func setupMiddleware(t *testing.T, subjectID string) (*bastion.Bastion, func()) {
	userService, teardown := setupService(t)

	// save a user to test
	u := user.User{Email: testUserEmail}
	err := u.SetPassword(testUserPassword)
	require.Nil(t, err)
	userService.Save(&u)

	if subjectID == "" {
		subjectID = u.ID.String()
	}

	middleware := getNewMiddleware(userService)

	app := bastion.New()
	app.APIRouter.Route("/", func(r chi.Router) {
		r.Use(setCtxMidd(subjectID))
		r.Use(middleware.LoggedUser)
		r.Get("/", handler)
	})

	return app, teardown
}

func TestLoggedUserMiddlewareOK(t *testing.T) {
	t.Parallel()

	app, teardown := setupMiddleware(t, "")
	defer teardown()

	e := bastion.Tester(t, app)
	e.GET("/").
		Expect().
		Status(http.StatusOK).
		Text().Equal(testUserEmail)
}

func TestLoggedUserMiddlewareUserNotFound(t *testing.T) {
	t.Parallel()

	// pass a new ulid to setupMiddleware to set another subjectID and force not found user
	app, teardown := setupMiddleware(t, kallax.NewULID().String())
	defer teardown()

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
