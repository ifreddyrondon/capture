package getting_test

import (
	"fmt"
	"testing"

	"github.com/ifreddyrondon/capture/pkg/domain"
	"github.com/ifreddyrondon/capture/pkg/getting"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"gopkg.in/src-d/go-kallax.v1"
)

type authorizationErr interface{ IsNotAuthorized() bool }

type mockStore struct {
	repo *domain.Repository
	err  error
}

func (m *mockStore) Get(kallax.ULID) (*domain.Repository, error) { return m.repo, m.err }

func TestServiceGetRepoOKWhenUserOwner(t *testing.T) {
	t.Parallel()

	userIDTxt := "0162eb39-a65e-04a1-7ad9-d663bb49a396"
	userID, err := kallax.NewULIDFromText(userIDTxt)
	assert.Nil(t, err)
	repoID := kallax.NewULID()
	store := &mockStore{repo: &domain.Repository{ID: repoID, Name: "test1", UserID: userID}}
	s := getting.NewService(store)

	u := &domain.User{ID: userID}
	repo, err := s.GetRepo(repoID, u)
	assert.Nil(t, err)
	assert.Equal(t, "test1", repo.Name)
}

func TestServiceGetRepoOKWhenPublic(t *testing.T) {
	t.Parallel()

	repoID := kallax.NewULID()
	store := &mockStore{repo: &domain.Repository{ID: repoID, Name: "test1", Visibility: domain.Public}}
	s := getting.NewService(store)

	userIDTxt := "0162eb39-a65e-04a1-7ad9-d663bb49a396"
	userID, err := kallax.NewULIDFromText(userIDTxt)
	assert.Nil(t, err)
	u := &domain.User{ID: userID}
	repo, err := s.GetRepo(repoID, u)
	assert.Nil(t, err)
	assert.Equal(t, "test1", repo.Name)
}

func TestServiceGetRepoErrWhenNoOwnerAndNoPublic(t *testing.T) {
	t.Parallel()

	repoID := kallax.NewULID()
	store := &mockStore{repo: &domain.Repository{ID: repoID, Name: "test1", Visibility: domain.Private}}
	s := getting.NewService(store)

	userIDTxt := "0162eb39-a65e-04a1-7ad9-d663bb49a396"
	userID, err := kallax.NewULIDFromText(userIDTxt)
	assert.Nil(t, err)
	u := &domain.User{ID: userID}
	_, err = s.GetRepo(repoID, u)
	assert.EqualError(t, err, fmt.Sprintf("user %v not authorized to get repo %v", userID, repoID))
	authErr, ok := errors.Cause(err).(authorizationErr)
	assert.True(t, ok)
	assert.True(t, authErr.IsNotAuthorized())
}

func TestServiceGetRepoErrGettingRepoFromStorage(t *testing.T) {
	t.Parallel()

	store := &mockStore{err: errors.New("test")}
	s := getting.NewService(store)

	userIDTxt := "0162eb39-a65e-04a1-7ad9-d663bb49a396"
	userID, err := kallax.NewULIDFromText(userIDTxt)
	assert.Nil(t, err)
	u := &domain.User{ID: userID}
	_, err = s.GetRepo(kallax.NewULID(), u)
	assert.EqualError(t, err, "could not get repo: test")
}
