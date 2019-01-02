package listing_test

import (
	"errors"
	"testing"

	listingBastion "github.com/ifreddyrondon/bastion/middleware/listing"
	"github.com/ifreddyrondon/bastion/middleware/listing/paging"
	"github.com/ifreddyrondon/capture/pkg/domain"
	"github.com/ifreddyrondon/capture/pkg/listing"
	"github.com/stretchr/testify/assert"
	"gopkg.in/src-d/go-kallax.v1"
)

type mockRepoStore struct {
	repos []domain.Repository
	err   error
}

func (m *mockRepoStore) List(*domain.Listing) ([]domain.Repository, error) { return m.repos, m.err }

func TestServiceGetUserReposOK(t *testing.T) {
	t.Parallel()

	store := &mockRepoStore{repos: []domain.Repository{
		{Name: "test1"},
		{Name: "test2"},
	}}
	s := listing.NewRepoService(store)

	u := &domain.User{ID: kallax.NewULID()}
	l := &listingBastion.Listing{
		Paging: paging.Paging{
			Limit:  50,
			Offset: 0,
		},
	}
	repos, err := s.GetUserRepos(u, l)
	assert.Nil(t, err)
	assert.NotNil(t, repos.Listing)
	assert.Equal(t, 2, len(repos.Results))
}

func TestServiceGetUserReposOKWhenEmpty(t *testing.T) {
	t.Parallel()

	store := &mockRepoStore{repos: nil}
	s := listing.NewRepoService(store)

	u := &domain.User{ID: kallax.NewULID()}
	l := &listingBastion.Listing{
		Paging: paging.Paging{
			Limit:  50,
			Offset: 0,
		},
	}
	repos, err := s.GetUserRepos(u, l)
	assert.Nil(t, err)
	assert.NotNil(t, repos.Listing)
	assert.Equal(t, 0, len(repos.Results))
}

func TestServiceGetUserReposErrWhenList(t *testing.T) {
	t.Parallel()

	store := &mockRepoStore{err: errors.New("test")}
	s := listing.NewRepoService(store)

	u := &domain.User{ID: kallax.NewULID()}
	l := &listingBastion.Listing{
		Paging: paging.Paging{
			Limit:  50,
			Offset: 0,
		},
	}
	_, err := s.GetUserRepos(u, l)
	assert.EqualError(t, err, "err getting user repos: test")
}

func TestServiceGetPublicReposOK(t *testing.T) {
	t.Parallel()

	store := &mockRepoStore{repos: []domain.Repository{
		{Name: "test1"},
		{Name: "test2"},
	}}
	s := listing.NewRepoService(store)
	l := &listingBastion.Listing{
		Paging: paging.Paging{
			Limit:  50,
			Offset: 0,
		},
	}
	repos, err := s.GetPublicRepos(l)
	assert.Nil(t, err)
	assert.NotNil(t, repos.Listing)
	assert.Equal(t, 2, len(repos.Results))
}

func TestServiceGetPublicReposErrWhenList(t *testing.T) {
	t.Parallel()

	store := &mockRepoStore{err: errors.New("test")}
	s := listing.NewRepoService(store)

	l := &listingBastion.Listing{
		Paging: paging.Paging{
			Limit:  50,
			Offset: 0,
		},
	}
	_, err := s.GetPublicRepos(l)
	assert.EqualError(t, err, "err getting public repos: test")
}
