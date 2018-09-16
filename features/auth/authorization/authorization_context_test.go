package authorization_test

import (
	"context"
	"testing"

	"github.com/ifreddyrondon/capture/features/auth/authorization"
	"github.com/stretchr/testify/assert"
)

func TestContextManagerGetCaptureOK(t *testing.T) {
	ctx := context.Background()

	ctx = authorization.WithSubjectID(ctx, "123")

	subj, err := authorization.GetSubjectID(ctx)
	assert.Nil(t, err)
	assert.Equal(t, subj, "123")
}

func TestContextManagerGetCaptureMissingCapture(t *testing.T) {
	ctx := context.Background()

	_, err := authorization.GetSubjectID(ctx)
	assert.EqualError(t, err, "subject id not found in context")
}
