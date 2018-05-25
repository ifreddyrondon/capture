package collection_test

import (
	"testing"

	"github.com/ifreddyrondon/gocapture/collection"

	"github.com/stretchr/testify/assert"
)

type MockRespository struct{}

func (r *MockRespository) Save(c *collection.Collection) error { return nil }

func setupServiceMockRepo(t *testing.T) *collection.REPOService {
	repo := &MockRespository{}
	return collection.NewService(repo)
}

func setupService(t *testing.T) (*collection.REPOService, func()) {
	repo, teardown := setupRepository(t)
	return collection.NewService(repo), teardown
}

func TestSaveSuccess(t *testing.T) {
	s := setupServiceMockRepo(t)

	c := collection.Collection{Name: "test"}
	err := s.Save(&c)
	assert.Nil(t, err)
	assert.NotEmpty(t, c.ID)
}
