package capture_test

import (
	"testing"

	"github.com/ifreddyrondon/gocapture/capture"
)

func setupService(t *testing.T) (*capture.StoreService, func()) {
	store, teardown := setupStore(t)
	return capture.NewService(store), teardown
}
