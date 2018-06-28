package listing_test

import (
	"context"
	"testing"

	"github.com/ifreddyrondon/capture/app/listing"
	"github.com/stretchr/testify/assert"
)

func TestContextManagerGetCaptureMissingCapture(t *testing.T) {
	ctxManager := listing.NewContextManager()
	ctx := context.Background()

	_, err := ctxManager.GetParams(ctx)
	assert.EqualError(t, err, "listing params not found in context")
}
