package repository_test

import (
	"testing"

	"github.com/ifreddyrondon/capture/features/repository"

	"github.com/stretchr/testify/assert"
)

type MockStore struct{}

func (r *MockStore) Save(c *repository.Repository) error { return nil }

func setupServiceMockRepo(t *testing.T) *repository.StoreService {
	store := &MockStore{}
	return repository.NewService(store)
}

func setupService(t *testing.T) (*repository.StoreService, func()) {
	store, teardown := setupStore(t)
	return repository.NewService(store), teardown
}

func TestSaveSuccess(t *testing.T) {
	s := setupServiceMockRepo(t)

	r := repository.Repository{Name: "test"}
	err := s.Save(&r)
	assert.Nil(t, err)
	assert.NotEmpty(t, r.ID)
}
