package middleware_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/go-chi/chi"
	"github.com/ifreddyrondon/bastion"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"gopkg.in/src-d/go-kallax.v1"

	"github.com/ifreddyrondon/capture/pkg/domain"
	"github.com/ifreddyrondon/capture/pkg/getting"
	"github.com/ifreddyrondon/capture/pkg/http/rest/middleware"
)

func setupRepoCtx(service getting.RepoService) *bastion.Bastion {
	app := bastion.New()
	app.Route("/{id}", func(r chi.Router) {
		r.Use(middleware.RepoCtx(service))
		r.Get("/", handler)
		r.Post("/", handler)
	})
	return app
}

type mockGettingRepoService struct {
	repo *domain.Repository
	err  error
}

func (m *mockGettingRepoService) Get(kallax.ULID) (*domain.Repository, error) {
	return m.repo, m.err
}

func TestRepoCtxSuccess(t *testing.T) {
	t.Parallel()

	app := setupRepoCtx(&mockGettingRepoService{})
	e := bastion.Tester(t, app)
	e.GET("/0167c8a5-d308-8692-809d-b1ad4a2d9562").
		Expect().
		Status(http.StatusOK)
}

func TestRepoCtxFailBadRequestGettingRepoByInvalidIDErr(t *testing.T) {
	t.Parallel()
	app := setupRepoCtx(&mockGettingRepoService{})
	e := bastion.Tester(t, app)

	response := map[string]interface{}{
		"status":  400.0,
		"error":   "Bad Request",
		"message": "invalid repository id",
	}

	e.GET("/a").
		Expect().
		Status(http.StatusBadRequest).
		JSON().Object().Equal(response)
}

func TestRepoCtxFailNotFoundGettingRepo(t *testing.T) {
	t.Parallel()
	app := setupRepoCtx(&mockGettingRepoService{err: notFound("test")})
	e := bastion.Tester(t, app)

	response := map[string]interface{}{
		"status":  404.0,
		"error":   "Not Found",
		"message": "not found repository",
	}

	e.GET("/0162eb39-a65e-04a1-7ad9-d663bb49a396").
		Expect().
		Status(http.StatusNotFound).
		JSON().Object().Equal(response)
}

func TestRepoCtxFailInternalServerErrGettingRepo(t *testing.T) {
	t.Parallel()
	app := setupRepoCtx(&mockGettingRepoService{err: errors.New("test")})
	e := bastion.Tester(t, app)

	response := map[string]interface{}{
		"status":  500.0,
		"error":   "Internal Server Error",
		"message": "looks like something went wrong",
	}

	e.GET("/0162eb39-a65e-04a1-7ad9-d663bb49a396").
		Expect().
		Status(http.StatusInternalServerError).
		JSON().Object().Equal(response)
}

func TestContextGetRepoOK(t *testing.T) {
	ctx := context.Background()
	repo := domain.Repository{Name: "test"}
	ctx = context.WithValue(ctx, middleware.RepoCtxKey, &repo)

	r, err := middleware.GetRepo(ctx)
	assert.Nil(t, err)
	assert.Equal(t, r.Name, r.Name)
}

func TestContextGetRepoMissingRepo(t *testing.T) {
	ctx := context.Background()
	_, err := middleware.GetRepo(ctx)
	assert.EqualError(t, err, "repo not found in context")
}

func TestContextGetRepoWhenWrongRepoValue(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, middleware.RepoCtxKey, "test")

	_, err := middleware.GetRepo(ctx)
	assert.EqualError(t, err, "repo value set incorrectly in context")
}
