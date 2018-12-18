package listing_test

import (
	"errors"
	"testing"

	listingBastion "github.com/ifreddyrondon/bastion/middleware/listing"
	"github.com/ifreddyrondon/bastion/middleware/listing/paging"
	"github.com/ifreddyrondon/capture/pkg"
	"github.com/ifreddyrondon/capture/pkg/domain"
	"github.com/ifreddyrondon/capture/pkg/listing"
	"github.com/stretchr/testify/assert"
)

type mockStore struct {
	repos []pkg.Repository
	err   error
}

func (m *mockStore) List(*domain.Listing) ([]pkg.Repository, error) { return m.repos, m.err }

func TestServiceGetUserReposOK(t *testing.T) {
	t.Parallel()

	store := &mockStore{repos: []pkg.Repository{
		{Name: "test1"},
		{Name: "test2"},
	}}
	s := listing.NewService(store)

	u := &pkg.User{ID: "0162eb39-a65e-04a1-7ad9-d663bb49a396"}
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

func TestServiceGetUserReposErrWhenList(t *testing.T) {
	t.Parallel()

	store := &mockStore{err: errors.New("test")}
	s := listing.NewService(store)

	u := &pkg.User{ID: "0162eb39-a65e-04a1-7ad9-d663bb49a396"}
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

	store := &mockStore{repos: []pkg.Repository{
		{Name: "test1"},
		{Name: "test2"},
	}}
	s := listing.NewService(store)
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

	store := &mockStore{err: errors.New("test")}
	s := listing.NewService(store)

	l := &listingBastion.Listing{
		Paging: paging.Paging{
			Limit:  50,
			Offset: 0,
		},
	}
	_, err := s.GetPublicRepos(l)
	assert.EqualError(t, err, "err getting public repos: test")
}
