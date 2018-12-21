package middleware_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/go-chi/chi"
	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/capture/pkg/authorizing"
	"github.com/ifreddyrondon/capture/pkg/domain"
	"github.com/ifreddyrondon/capture/pkg/http/rest/middleware"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"gopkg.in/src-d/go-kallax.v1"
)

func setupAuthorizing(service authorizing.Service) *bastion.Bastion {
	app := bastion.New()
	app.APIRouter.Route("/", func(r chi.Router) {
		r.Use(middleware.AuthorizeReq(service))
		r.Get("/", handler)
		r.Post("/", handler)
	})
	return app
}

type mockAuthorizingService struct {
	usr *domain.User
	err error
}

func (m *mockAuthorizingService) AuthorizeRequest(*http.Request) (*domain.User, error) {
	return m.usr, m.err
}

func TestAuthorizingSuccess(t *testing.T) {
	t.Parallel()

	app := setupAuthorizing(&mockAuthorizingService{})
	e := bastion.Tester(t, app)
	e.GET("/").WithHeader("Authorization", fmt.Sprintf("Bearer %v", "test")).
		Expect().
		Status(http.StatusOK)
}

func TestAuthorizingInvalidID(t *testing.T) {
	t.Parallel()

	app := setupAuthorizing(&mockAuthorizingService{err: invalidErr("test")})
	response := map[string]interface{}{
		"status":  400.0,
		"error":   "Bad Request",
		"message": "invalid user id format",
	}

	e := bastion.Tester(t, app)
	e.GET("/").WithHeader("Authorization", "Bearer Test").
		Expect().
		Status(http.StatusBadRequest).
		JSON().Object().Equal(response)
}

func TestAuthorizingNotAuthorized(t *testing.T) {
	t.Parallel()

	app := setupAuthorizing(&mockAuthorizingService{err: notAllowedErr("you don’t have permission to access this resource")})
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

func TestAuthorizingNotFound(t *testing.T) {
	t.Parallel()

	app := setupAuthorizing(&mockAuthorizingService{err: notFound("test")})
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

	app := setupAuthorizing(&mockAuthorizingService{err: errors.New("test")})
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
	u := domain.User{ID: kallax.NewULID(), Email: "test@example.com"}
	ctx = context.WithValue(ctx, middleware.UserCtxKey, &u)

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
