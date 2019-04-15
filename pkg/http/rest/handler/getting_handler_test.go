package handler_test

import (
	"net/http"
	"testing"

	"github.com/ifreddyrondon/bastion"

	"github.com/ifreddyrondon/capture/pkg/http/rest/handler"
)

func setupGettingRepoHandler(m func(http.Handler) http.Handler) *bastion.Bastion {
	app := bastion.New()
	app.Use(m)
	app.Get("/{id}", handler.GettingRepo())
	return app
}

func TestGettingRepoSuccess(t *testing.T) {
	t.Parallel()

	app := setupGettingRepoHandler(withRepoMiddle(defaultRepo))

	e := bastion.Tester(t, app)
	e.GET("/0167c8a5-d308-8692-809d-b1ad4a2d9562").
		Expect().
		JSON().Object().
		ContainsKey("name").ValueEqual("name", "test public").
		ContainsKey("visibility").ValueEqual("visibility", "public")
}

func TestGettingRepoInternalServer(t *testing.T) {
	t.Parallel()

	app := setupGettingRepoHandler(withRepoMiddle(nil))

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

func setupGettingCaptureHandler(m func(http.Handler) http.Handler) *bastion.Bastion {
	app := bastion.New()
	app.Use(m)
	app.Get("/{id}", handler.GettingCapture())
	return app
}

func TestGettingCaptureSuccess(t *testing.T) {
	t.Parallel()

	app := setupGettingCaptureHandler(withCaptureMiddle(defaultCapture))

	e := bastion.Tester(t, app)
	e.GET("/0167c8a5-d308-8692-809d-b1ad4a2d9562").
		Expect().
		JSON().Object().
		ContainsKey("id")
}

func TestGettingCaptureInternalServer(t *testing.T) {
	t.Parallel()

	app := setupGettingCaptureHandler(withCaptureMiddle(nil))

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
