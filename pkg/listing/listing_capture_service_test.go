package listing_test

import (
	"fmt"
	"testing"

	listingBastion "github.com/ifreddyrondon/bastion/middleware/listing"
	"github.com/ifreddyrondon/bastion/middleware/listing/paging"
	"github.com/ifreddyrondon/capture/pkg/domain"
	"github.com/ifreddyrondon/capture/pkg/listing"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"gopkg.in/src-d/go-kallax.v1"
)

type mockCaptureStore struct {
	captures []domain.Capture
	err      error
}

func (m *mockCaptureStore) List(*domain.Listing) ([]domain.Capture, error) { return m.captures, m.err }

func TestCaptureServiceListRepoCapturesOKWhenRepoBelongToUser(t *testing.T) {
	t.Parallel()

	store := &mockCaptureStore{captures: []domain.Capture{
		{ID: kallax.NewULID()}, {ID: kallax.NewULID()},
	}}
	s := listing.NewCaptureService(store)
	userID := kallax.NewULID()

	u := &domain.User{ID: userID}
	r := &domain.Repository{Name: "test", ID: kallax.NewULID(), UserID: userID, Visibility: domain.Private}
	l := &listingBastion.Listing{
		Paging: paging.Paging{
			Limit:  50,
			Offset: 0,
		},
	}
	captures, err := s.ListRepoCaptures(u, r, l)
	assert.Nil(t, err)
	assert.NotNil(t, captures.Listing)
	assert.Equal(t, 2, len(captures.Results))
}

func TestCaptureServiceListRepoCapturesOKWhenRepoPublic(t *testing.T) {
	t.Parallel()

	store := &mockCaptureStore{captures: []domain.Capture{
		{ID: kallax.NewULID()}, {ID: kallax.NewULID()},
	}}
	s := listing.NewCaptureService(store)
	u := &domain.User{ID: kallax.NewULID()}
	r := &domain.Repository{Name: "test", ID: kallax.NewULID(), UserID: kallax.NewULID(), Visibility: domain.Public}
	l := &listingBastion.Listing{
		Paging: paging.Paging{
			Limit:  50,
			Offset: 0,
		},
	}
	captures, err := s.ListRepoCaptures(u, r, l)
	assert.Nil(t, err)
	assert.NotNil(t, captures.Listing)
	assert.Equal(t, 2, len(captures.Results))
}

func TestCaptureServiceListRepoCapturesOKWhenEmpty(t *testing.T) {
	t.Parallel()

	store := &mockCaptureStore{captures: nil}
	s := listing.NewCaptureService(store)
	u := &domain.User{ID: kallax.NewULID()}
	r := &domain.Repository{Name: "test", ID: kallax.NewULID(), UserID: kallax.NewULID(), Visibility: domain.Public}
	l := &listingBastion.Listing{
		Paging: paging.Paging{
			Limit:  50,
			Offset: 0,
		},
	}
	captures, err := s.ListRepoCaptures(u, r, l)
	assert.Nil(t, err)
	assert.NotNil(t, captures.Listing)
	assert.Equal(t, 0, len(captures.Results))
}

type notAuthorizedErr interface {
	IsNotAuthorized() bool
}

func TestCaptureServiceListRepoCapturesErrWhenNotAuth(t *testing.T) {
	t.Parallel()

	store := &mockCaptureStore{captures: []domain.Capture{
		{ID: kallax.NewULID()}, {ID: kallax.NewULID()},
	}}
	s := listing.NewCaptureService(store)
	userID, repoID := kallax.NewULID(), kallax.NewULID()

	u := &domain.User{ID: userID}
	r := &domain.Repository{Name: "test", ID: repoID, Visibility: domain.Private}
	l := &listingBastion.Listing{
		Paging: paging.Paging{
			Limit:  50,
			Offset: 0,
		},
	}
	_, err := s.ListRepoCaptures(u, r, l)
	assert.EqualError(t, err, fmt.Sprintf("user %v not authorized to get captures from repo %v", userID, repoID))
	authErr, ok := errors.Cause(err).(notAuthorizedErr)
	assert.True(t, ok)
	assert.True(t, authErr.IsNotAuthorized())
}

func TestCaptureServiceListRepoCapturesErrWhenList(t *testing.T) {
	t.Parallel()

	store := &mockCaptureStore{err: errors.New("test")}
	s := listing.NewCaptureService(store)
	userID := kallax.NewULID()

	u := &domain.User{ID: userID}
	r := &domain.Repository{Name: "test", ID: kallax.NewULID(), UserID: userID, Visibility: domain.Private}
	l := &listingBastion.Listing{
		Paging: paging.Paging{
			Limit:  50,
			Offset: 0,
		},
	}
	_, err := s.ListRepoCaptures(u, r, l)
	assert.EqualError(t, err, "err getting repo captures: test")
}
