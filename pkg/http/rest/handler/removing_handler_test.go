package handler_test

import (
	"net/http"
	"testing"

	"github.com/ifreddyrondon/capture/pkg/http/rest/handler"

	"github.com/ifreddyrondon/bastion"
	"github.com/pkg/errors"

	"github.com/ifreddyrondon/capture/pkg/domain"
	"github.com/ifreddyrondon/capture/pkg/removing"
)

type removingCaptureServiceMock struct {
	err error
}

func (m *removingCaptureServiceMock) Remove(*domain.Capture) error { return m.err }

func setupRemovingCaptureHandler(s removing.CaptureService, m func(http.Handler) http.Handler) *bastion.Bastion {
	app := bastion.New()
	app.Use(m)
	app.Get("/", handler.RemovingCapture(s))
	return app
}

func TestRemovingCaptureSuccess(t *testing.T) {
	t.Parallel()

	app := setupRemovingCaptureHandler(&removingCaptureServiceMock{}, withCaptureMiddle(defaultCapture))

	e := bastion.Tester(t, app)
	e.GET("/").
		Expect().
		JSON().Object().
		ContainsKey("id")
}

func TestRemovingCaptureFailsGettingCapture(t *testing.T) {
	t.Parallel()

	app := setupRemovingCaptureHandler(&removingCaptureServiceMock{}, withCaptureMiddle(nil))

	response := map[string]interface{}{
		"status":  500.0,
		"error":   "Internal Server Error",
		"message": "looks like something went wrong",
	}

	e := bastion.Tester(t, app)
	e.GET("/").Expect().
		Status(http.StatusInternalServerError).
		JSON().Object().Equal(response)
}

func TestRemovingCaptureFailsRemoving(t *testing.T) {
	t.Parallel()
	s := &removingCaptureServiceMock{err: errors.New("test")}
	app := setupRemovingCaptureHandler(s, withCaptureMiddle(defaultCapture))

	response := map[string]interface{}{
		"status":  500.0,
		"error":   "Internal Server Error",
		"message": "looks like something went wrong",
	}

	e := bastion.Tester(t, app)
	e.GET("/").Expect().
		Status(http.StatusInternalServerError).
		JSON().Object().Equal(response)
}
