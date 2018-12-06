package middleware_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/go-chi/chi"
	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/capture/pkg"
	"github.com/ifreddyrondon/capture/pkg/authorizing"
	"github.com/ifreddyrondon/capture/pkg/http/rest/middleware"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"gopkg.in/src-d/go-kallax.v1"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "ok")
}

func setup(service authorizing.Service) *bastion.Bastion {
	app := bastion.New()
	app.APIRouter.Route("/", func(r chi.Router) {
		r.Use(middleware.AuthorizeReq(service))
		r.Get("/", handler)
		r.Post("/", handler)
	})
	return app
}

type mockService struct {
	usr *pkg.User
	err error
}

func (m *mockService) AuthorizeRequest(*http.Request) (*pkg.User, error) { return m.usr, m.err }

func TestAuthorizingSuccess(t *testing.T) {
	t.Parallel()

	app := setup(&mockService{})
	e := bastion.Tester(t, app)
	e.GET("/").WithHeader("Authorization", fmt.Sprintf("Bearer %v", "test")).
		Expect().
		Status(http.StatusOK)
}

type invalidCredentialErr string

func (i invalidCredentialErr) Error() string         { return fmt.Sprintf(string(i)) }
func (i invalidCredentialErr) IsNotAuthorized() bool { return true }

func TestAuthorizingNotAuthorized(t *testing.T) {
	t.Parallel()

	app := setup(&mockService{err: invalidCredentialErr("you don’t have permission to access this resource")})
	response := map[string]interface{}{
		"status":  403.0,
		"error":   "Forbidden",
		"message": "you don’t have permission to access this resource",
	}

	e := bastion.Tester(t, app)
	e.GET("/").WithHeader("Authorization", "Bearer Test").
		Expect().
		Status(http.StatusForbidden).
		JSON().Object().Equal(response)
}

type userNotFound string

func (u userNotFound) Error() string  { return fmt.Sprintf(string(u)) }
func (u userNotFound) NotFound() bool { return true }

func TestAuthorizingNotFound(t *testing.T) {
	t.Parallel()

	app := setup(&mockService{err: userNotFound("test")})
	response := map[string]interface{}{
		"status":  404.0,
		"error":   "Not Found",
		"message": "test",
	}

	e := bastion.Tester(t, app)
	e.GET("/").WithHeader("Authorization", "Bearer Test").
		Expect().
		Status(http.StatusNotFound).
		JSON().Object().Equal(response)
}

func TestAuthorizingInternalErr(t *testing.T) {
	t.Parallel()

	app := setup(&mockService{err: errors.New("test")})
	response := map[string]interface{}{
		"status":  500.0,
		"error":   "Internal Server Error",
		"message": "looks like something went wrong",
	}

	e := bastion.Tester(t, app)
	e.GET("/").WithHeader("Authorization", "Bearer Test").
		Expect().
		Status(http.StatusInternalServerError).
		JSON().Object().Equal(response)
}

func TestContextGetUserOK(t *testing.T) {
	ctx := context.Background()
	u := pkg.User{ID: kallax.NewULID(), Email: "test@example.com"}
	ctx = middleware.WithUser(ctx, &u)

	u2, err := middleware.GetUser(ctx)
	assert.Nil(t, err)
	assert.Equal(t, u.ID, u2.ID)
	assert.Equal(t, u.Email, u2.Email)
}

func TestContextGetUserMissingUser(t *testing.T) {
	ctx := context.Background()
	_, err := middleware.GetUser(ctx)
	assert.EqualError(t, err, "user not found in context")
}

func TestContextGetUserIDMissingUser(t *testing.T) {
	ctx := context.Background()

	_, err := middleware.GetUser(ctx)
	assert.EqualError(t, err, "user not found in context")
}
