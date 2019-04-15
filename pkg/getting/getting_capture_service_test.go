package getting_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/src-d/go-kallax.v1"

	"github.com/ifreddyrondon/capture/pkg/domain"
	"github.com/ifreddyrondon/capture/pkg/getting"
)

type mockCaptureStore struct {
	capture *domain.Capture
	err     error
}

func (m *mockCaptureStore) Get(captureID, repoID kallax.ULID) (*domain.Capture, error) {
	return m.capture, m.err
}

func TestServiceGetCaptureOK(t *testing.T) {
	t.Parallel()

	captID := kallax.NewULID()
	store := &mockCaptureStore{capture: &domain.Capture{ID: captID}}
	s := getting.NewCaptureService(store)
	r := &domain.Repository{ID: kallax.NewULID(), Name: "test1"}

	repo, err := s.Get(captID, r)
	assert.Nil(t, err)
	assert.Equal(t, captID, repo.ID)
}

func TestServiceGetCaptureErrorGettingTheCapture(t *testing.T) {
	t.Parallel()

	store := &mockCaptureStore{err: errors.New("test")}
	s := getting.NewCaptureService(store)
	r := &domain.Repository{ID: kallax.NewULID(), Name: "test1"}

	_, err := s.Get(kallax.NewULID(), r)
	assert.EqualError(t, err, "could not get capture: test")
}
