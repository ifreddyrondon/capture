package rest_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/capture/pkg"
	"github.com/ifreddyrondon/capture/pkg/http/rest"
	"github.com/ifreddyrondon/capture/pkg/http/rest/middleware"
)

func repoCtxtMiddlewareOK(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		repo := &pkg.Repository{Name: "test public", Visibility: "public"}

		ctx := context.WithValue(r.Context(), middleware.RepoCtxKey, repo)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func repoCtxtMiddlewareBAD(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), middleware.RepoCtxKey, "bad listing")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func setupGettingHandler(m func(http.Handler) http.Handler) *bastion.Bastion {
	app := bastion.New()
	app.APIRouter.Use(m)
	app.APIRouter.Get("/{id}", rest.GettingRepo())
	return app
}

func TestGettingRepoSuccess(t *testing.T) {
	t.Parallel()

	app := setupGettingHandler(repoCtxtMiddlewareOK)

	e := bastion.Tester(t, app)
	e.GET("/0167c8a5-d308-8692-809d-b1ad4a2d9562").
		Expect().
		JSON().Object().
		ContainsKey("name").ValueEqual("name", "test public").
		ContainsKey("visibility").ValueEqual("visibility", "public")
}

func TestLGettingRepoInternalServer(t *testing.T) {
	t.Parallel()

	app := setupGettingHandler(repoCtxtMiddlewareBAD)

	response := map[string]interface{}{
		"status":  500.0,
		"error":   "Internal Server Error",
		"message": "looks like something went wrong",
	}

	e := bastion.Tester(t, app)
	e.GET("/0167c8a5-d308-8692-809d-b1ad4a2d9562").Expect().
		Status(http.StatusInternalServerError).
		JSON().Object().Equal(response)
}
