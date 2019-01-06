package removing_test

import (
	"testing"
	"time"

	"github.com/ifreddyrondon/capture/pkg/domain"
	"github.com/ifreddyrondon/capture/pkg/removing"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"gopkg.in/src-d/go-kallax.v1"
)

type mockCaptureStore struct {
	err error
}

func (m *mockCaptureStore) Save(*domain.Capture) error {
	return m.err
}

func TestServiceRemoveCaptureOK(t *testing.T) {
	t.Parallel()

	store := &mockCaptureStore{}
	s := removing.NewCaptureService(store)
	capt := &domain.Capture{ID: kallax.NewULID()}

	timeBeforeDelete := time.Now()
	err := s.Remove(capt)
	assert.Nil(t, err)
	assert.NotNil(t, capt.DeletedAt)
	assert.True(t, capt.DeletedAt.After(timeBeforeDelete))
}

func TestServiceRemoveCaptureFailsWhenSave(t *testing.T) {
	t.Parallel()

	store := &mockCaptureStore{err: errors.New("test")}
	s := removing.NewCaptureService(store)
	capt := &domain.Capture{ID: kallax.NewULID()}

	err := s.Remove(capt)
	assert.EqualError(t, err, "could not remove capture: test")
}
