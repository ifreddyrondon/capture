package capture_test

import (
	"testing"

	"github.com/ifreddyrondon/gocapture/capture"
)

func setupService(t *testing.T) (*capture.REPOService, func()) {
	repo, teardown := setupRepository(t)
	return capture.NewService(repo), teardown
}
