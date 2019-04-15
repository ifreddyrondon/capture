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

func setupCaptureCtx(service getting.CaptureService, getRepo func(http.Handler) http.Handler) *bastion.Bastion {
	app := bastion.New()
	app.Route("/{captureId}", func(r chi.Router) {
		r.Use(getRepo)
		r.Use(middleware.CaptureCtx(service))
		r.Get("/", handler)
		r.Post("/", handler)
	})
	return app
}

type mockGettingCaptureService struct {
	capt *domain.Capture
	err  error
}

func (m *mockGettingCaptureService) Get(kallax.ULID, *domain.Repository) (*domain.Capture, error) {
	return m.capt, m.err
}

func TestCaptureCtxSuccess(t *testing.T) {
	t.Parallel()

	app := setupCaptureCtx(&mockGettingCaptureService{}, withRepoMiddle(defaultRepo))
	e := bastion.Tester(t, app)
	e.GET("/0167c8a5-d308-8692-809d-b1ad4a2d9562").
		Expect().
		Status(http.StatusOK)
}

func TestCaptureCtxFailInternalErrorGettingRepo(t *testing.T) {
	t.Parallel()
	s := &mockGettingCaptureService{}
	app := setupCaptureCtx(s, withRepoMiddle(nil))
	e := bastion.Tester(t, app)

	response := map[string]interface{}{
		"status":  500.0,
		"error":   "Internal Server Error",
		"message": "looks like something went wrong",
	}

	e.GET("/0167c8a5-d308-8692-809d-b1ad4a2d9562").
		Expect().
		Status(http.StatusInternalServerError).
		JSON().Object().Equal(response)
}

func TestCaptureCtxFailBadRequestGettingCaptureByInvalidIDErr(t *testing.T) {
	t.Parallel()
	app := setupCaptureCtx(&mockGettingCaptureService{}, withRepoMiddle(defaultRepo))
	e := bastion.Tester(t, app)

	response := map[string]interface{}{
		"status":  400.0,
		"error":   "Bad Request",
		"message": "invalid capture id",
	}

	e.GET("/a").
		Expect().
		Status(http.StatusBadRequest).
		JSON().Object().Equal(response)
}

func TestCaptureCtxFailNotFoundGettingCapture(t *testing.T) {
	t.Parallel()
	s := &mockGettingCaptureService{err: notFound("test")}
	app := setupCaptureCtx(s, withRepoMiddle(defaultRepo))
	e := bastion.Tester(t, app)

	response := map[string]interface{}{
		"status":  404.0,
		"error":   "Not Found",
		"message": "not found capture",
	}

	e.GET("/0162eb39-a65e-04a1-7ad9-d663bb49a396").
		Expect().
		Status(http.StatusNotFound).
		JSON().Object().Equal(response)
}

func TestCaptureCtxFailInternalServerErrGettingCapture(t *testing.T) {
	t.Parallel()
	s := &mockGettingCaptureService{err: errors.New("test")}
	app := setupCaptureCtx(s, withRepoMiddle(defaultRepo))
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

func TestContextGetCaptureOK(t *testing.T) {
	captID := kallax.NewULID()
	ctx := context.Background()
	capt := domain.Capture{ID: captID}
	ctx = context.WithValue(ctx, middleware.CaptureCtxKey, &capt)

	c, err := middleware.GetCapture(ctx)
	assert.Nil(t, err)
	assert.Equal(t, captID, c.ID)
}

func TestContextGetCaptureMissingCapture(t *testing.T) {
	ctx := context.Background()
	_, err := middleware.GetCapture(ctx)
	assert.EqualError(t, err, "capture not found in context")
}

func TestContextGetCaptureWhenWrongCaptureValue(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, middleware.CaptureCtxKey, "test")

	_, err := middleware.GetCapture(ctx)
	assert.EqualError(t, err, "capture value set incorrectly in context")
}
