package getting_test

import (
	"testing"

	"github.com/ifreddyrondon/capture/pkg/domain"
	"github.com/ifreddyrondon/capture/pkg/getting"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"gopkg.in/src-d/go-kallax.v1"
)

type mockRepoStore struct {
	repo *domain.Repository
	err  error
}

func (m *mockRepoStore) Get(kallax.ULID) (*domain.Repository, error) { return m.repo, m.err }

func TestServiceGetRepoOK(t *testing.T) {
	t.Parallel()

	repoID := kallax.NewULID()
	store := &mockRepoStore{repo: &domain.Repository{ID: repoID, Name: "test1"}}
	s := getting.NewRepoService(store)
	repo, err := s.Get(repoID)
	assert.Nil(t, err)
	assert.Equal(t, "test1", repo.Name)
}

func TestServiceGetRepoErrGettingRepoFromStorage(t *testing.T) {
	t.Parallel()

	store := &mockRepoStore{err: errors.New("test")}
	s := getting.NewRepoService(store)
	_, err := s.Get(kallax.NewULID())
	assert.EqualError(t, err, "could not get repo: test")
}
