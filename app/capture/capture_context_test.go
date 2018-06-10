package capture_test

import (
	"context"
	"testing"

	"github.com/ifreddyrondon/capture/app/capture"

	"github.com/stretchr/testify/assert"
	kallax "gopkg.in/src-d/go-kallax.v1"
)

func TestContextManagerGetCaptureOK(t *testing.T) {
	ctxManager := capture.NewContextManager()
	ctx := context.Background()

	cap := capture.Capture{ID: kallax.NewULID()}
	ctx = ctxManager.WithCapture(ctx, &cap)

	cap2, err := ctxManager.GetCapture(ctx)
	assert.Nil(t, err)
	assert.Equal(t, cap.ID, cap2.ID)
}

func TestContextManagerGetCaptureMissingCapture(t *testing.T) {
	ctxManager := capture.NewContextManager()
	ctx := context.Background()

	_, err := ctxManager.GetCapture(ctx)
	assert.EqualError(t, err, "capture not found in context")
}
