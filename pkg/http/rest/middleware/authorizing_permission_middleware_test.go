package middleware_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/go-chi/chi"
	"github.com/ifreddyrondon/bastion"
	"gopkg.in/src-d/go-kallax.v1"

	"github.com/ifreddyrondon/capture/pkg/domain"
	"github.com/ifreddyrondon/capture/pkg/http/rest/middleware"
)

func setupRepoOwner(auth, repoCtx func(http.Handler) http.Handler) *bastion.Bastion {
	app := bastion.New()
	app.Route("/", func(r chi.Router) {
		r.Use(auth)
		r.Use(repoCtx)
		r.Use(middleware.RepoOwner())
		r.Get("/", handler)
	})
	return app
}

func TestRepoOwnerSuccessWhenRepoBelongToUser(t *testing.T) {
	t.Parallel()
	repo := &domain.Repository{Name: "test", Visibility: domain.Private, UserID: defaultUserID}
	app := setupRepoOwner(withUserMiddle(defaultUser), withRepoMiddle(repo))
	e := bastion.Tester(t, app)
	e.GET("/").
		Expect().
		Status(http.StatusOK)
}

func TestRepoOwnerFailsBecauseDontBelong2User(t *testing.T) {
	t.Parallel()
	repoID := kallax.NewULID()
	repo := &domain.Repository{ID: repoID, Name: "test", Visibility: domain.Private, UserID: kallax.NewULID()}
	app := setupRepoOwner(withUserMiddle(defaultUser), withRepoMiddle(repo))
	e := bastion.Tester(t, app)

	response := map[string]interface{}{
		"status":  403.0,
		"error":   "Forbidden",
		"message": fmt.Sprintf("You don't have permission to access repository %v", repoID),
	}

	e.GET("/").
		Expect().
		Status(http.StatusForbidden).
		JSON().Object().Equal(response)
}

func TestRepoOwnerFailsBecauseGetUserMiddleware(t *testing.T) {
	t.Parallel()
	app := setupRepoOwner(withUserMiddle(nil), withRepoMiddle(nil))
	e := bastion.Tester(t, app)

	response := map[string]interface{}{
		"status":  500.0,
		"error":   "Internal Server Error",
		"message": "looks like something went wrong",
	}

	e.GET("/").
		Expect().
		Status(http.StatusInternalServerError).
		JSON().Object().Equal(response)
}

func TestRepoOwnerFailsBecauseGetRepoMiddleware(t *testing.T) {
	t.Parallel()
	app := setupRepoOwner(withUserMiddle(defaultUser), withRepoMiddle(nil))
	e := bastion.Tester(t, app)

	response := map[string]interface{}{
		"status":  500.0,
		"error":   "Internal Server Error",
		"message": "looks like something went wrong",
	}

	e.GET("/").
		Expect().
		Status(http.StatusInternalServerError).
		JSON().Object().Equal(response)
}

func setupRepoOwnerOrPublic(auth, repoCtx func(http.Handler) http.Handler) *bastion.Bastion {
	app := bastion.New()
	app.Route("/", func(r chi.Router) {
		r.Use(auth)
		r.Use(repoCtx)
		r.Use(middleware.RepoOwnerOrPublic())
		r.Get("/", handler)
	})
	return app
}

func TestRepoOwnerOrPublicSuccessWhenPublicRepo(t *testing.T) {
	t.Parallel()
	repo := &domain.Repository{Name: "test", Visibility: domain.Public}
	app := setupRepoOwnerOrPublic(withUserMiddle(defaultUser), withRepoMiddle(repo))
	e := bastion.Tester(t, app)
	e.GET("/").
		Expect().
		Status(http.StatusOK)
}

func TestRepoOwnerOrPublicSuccessWhenRepoBelongToUser(t *testing.T) {
	t.Parallel()
	repo := &domain.Repository{Name: "test", Visibility: domain.Private, UserID: defaultUserID}
	app := setupRepoOwnerOrPublic(withUserMiddle(defaultUser), withRepoMiddle(repo))
	e := bastion.Tester(t, app)
	e.GET("/").
		Expect().
		Status(http.StatusOK)
}

func TestRepoOwnerFailsBecauseDontBelong2UserAndPrivate(t *testing.T) {
	t.Parallel()
	repoID := kallax.NewULID()
	repo := &domain.Repository{ID: repoID, Name: "test", Visibility: domain.Private, UserID: kallax.NewULID()}
	app := setupRepoOwnerOrPublic(withUserMiddle(defaultUser), withRepoMiddle(repo))
	e := bastion.Tester(t, app)

	response := map[string]interface{}{
		"status":  403.0,
		"error":   "Forbidden",
		"message": fmt.Sprintf("You don't have permission to access repository %v", repoID),
	}

	e.GET("/").
		Expect().
		Status(http.StatusForbidden).
		JSON().Object().Equal(response)
}

func TestRepoOwnerOrPublicFailsBecauseGetUserMiddleware(t *testing.T) {
	t.Parallel()
	app := setupRepoOwnerOrPublic(withUserMiddle(nil), withRepoMiddle(nil))
	e := bastion.Tester(t, app)

	response := map[string]interface{}{
		"status":  500.0,
		"error":   "Internal Server Error",
		"message": "looks like something went wrong",
	}

	e.GET("/").
		Expect().
		Status(http.StatusInternalServerError).
		JSON().Object().Equal(response)
}

func TestRepoOwnerOrPublicFailsBecauseGetRepoMiddleware(t *testing.T) {
	t.Parallel()
	app := setupRepoOwnerOrPublic(withUserMiddle(defaultUser), withRepoMiddle(nil))
	e := bastion.Tester(t, app)

	response := map[string]interface{}{
		"status":  500.0,
		"error":   "Internal Server Error",
		"message": "looks like something went wrong",
	}

	e.GET("/").
		Expect().
		Status(http.StatusInternalServerError).
		JSON().Object().Equal(response)
}
