package listing_test

import (
	"context"
	"testing"

	"github.com/ifreddyrondon/capture/app/listing"
	"github.com/stretchr/testify/assert"
)

func TestContextManagerGetListingMissingInstance(t *testing.T) {
	ctxManager := listing.NewContextManager()
	ctx := context.Background()

	_, err := ctxManager.GetListing(ctx)
	assert.EqualError(t, err, "listing not found in context")
}
