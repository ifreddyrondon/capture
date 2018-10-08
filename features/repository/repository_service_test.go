package repository_test

import (
	"testing"

	"github.com/ifreddyrondon/capture/features"
	"github.com/ifreddyrondon/capture/features/repository"

	"github.com/stretchr/testify/assert"
)

func setupService(t *testing.T) (*repository.StoreService, func()) {
	store, teardown := setupStore(t)
	return repository.NewService(store), teardown
}

type mockStore struct{ err error }

func (r *mockStore) Save(c *features.Repository) error { return r.err }

func TestSaveSuccess(t *testing.T) {
	s := mockStore{nil}

	r := features.Repository{Name: "test"}
	err := s.Save(&r)
	assert.Nil(t, err)
	assert.NotEmpty(t, r.ID)
}
