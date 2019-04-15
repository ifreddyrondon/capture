package listing_test

import (
	"testing"

	listingBastion "github.com/ifreddyrondon/bastion/middleware/listing"
	"github.com/ifreddyrondon/bastion/middleware/listing/paging"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"gopkg.in/src-d/go-kallax.v1"

	"github.com/ifreddyrondon/capture/pkg/domain"
	"github.com/ifreddyrondon/capture/pkg/listing"
)

type mockCaptureStore struct {
	captures []domain.Capture
	count    int64
	err      error
}

func (m *mockCaptureStore) List(*domain.Listing) ([]domain.Capture, int64, error) {
	return m.captures, m.count, m.err
}

func TestCaptureServiceListRepoCapturesOK(t *testing.T) {
	t.Parallel()

	store := &mockCaptureStore{captures: []domain.Capture{
		{ID: kallax.NewULID()}, {ID: kallax.NewULID()},
	}}
	s := listing.NewCaptureService(store)
	userID := kallax.NewULID()

	r := &domain.Repository{Name: "test", ID: kallax.NewULID(), UserID: userID, Visibility: domain.Private}
	l := &listingBastion.Listing{
		Paging: paging.Paging{
			Limit:  50,
			Offset: 0,
		},
	}
	captures, err := s.ListRepoCaptures(r, l)
	assert.Nil(t, err)
	assert.NotNil(t, captures.Listing)
	assert.Equal(t, 2, len(captures.Results))
}

func TestCaptureServiceListRepoCapturesOKWhenEmpty(t *testing.T) {
	t.Parallel()

	store := &mockCaptureStore{captures: nil}
	s := listing.NewCaptureService(store)
	r := &domain.Repository{Name: "test", ID: kallax.NewULID(), UserID: kallax.NewULID(), Visibility: domain.Public}
	l := &listingBastion.Listing{
		Paging: paging.Paging{
			Limit:  50,
			Offset: 0,
		},
	}
	captures, err := s.ListRepoCaptures(r, l)
	assert.Nil(t, err)
	assert.NotNil(t, captures.Listing)
	assert.Equal(t, 0, len(captures.Results))
}

func TestCaptureServiceListRepoCapturesErrWhenList(t *testing.T) {
	t.Parallel()

	store := &mockCaptureStore{err: errors.New("test")}
	s := listing.NewCaptureService(store)
	userID := kallax.NewULID()

	r := &domain.Repository{Name: "test", ID: kallax.NewULID(), UserID: userID, Visibility: domain.Private}
	l := &listingBastion.Listing{
		Paging: paging.Paging{
			Limit:  50,
			Offset: 0,
		},
	}
	_, err := s.ListRepoCaptures(r, l)
	assert.EqualError(t, err, "err getting repo captures: test")
}
