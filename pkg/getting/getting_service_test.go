package getting_test

import (
	"fmt"
	"testing"

	"github.com/ifreddyrondon/capture/pkg"
	"github.com/ifreddyrondon/capture/pkg/getting"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"gopkg.in/src-d/go-kallax.v1"
)

type authorizationErr interface{ IsNotAuthorized() bool }

type mockStore struct {
	repo *pkg.Repository
	err  error
}

func (m *mockStore) Get(string) (*pkg.Repository, error) { return m.repo, m.err }

func TestServiceGetRepoOKWhenUserOwner(t *testing.T) {
	t.Parallel()

	userIDTxt := "0162eb39-a65e-04a1-7ad9-d663bb49a396"
	userID, err := kallax.NewULIDFromText(userIDTxt)
	assert.Nil(t, err)
	repoID := kallax.NewULID()
	store := &mockStore{repo: &pkg.Repository{ID: repoID, Name: "test1", UserID: userID}}
	s := getting.NewService(store)

	u := &pkg.User{ID: userIDTxt}
	repo, err := s.GetRepo(repoID.String(), u)
	assert.Nil(t, err)
	assert.Equal(t, "test1", repo.Name)
}

func TestServiceGetRepoOKWhenPublic(t *testing.T) {
	t.Parallel()

	repoID := kallax.NewULID()
	store := &mockStore{repo: &pkg.Repository{ID: repoID, Name: "test1", Visibility: pkg.Public}}
	s := getting.NewService(store)

	u := &pkg.User{ID: "0162eb39-a65e-04a1-7ad9-d663bb49a396"}
	repo, err := s.GetRepo(repoID.String(), u)
	assert.Nil(t, err)
	assert.Equal(t, "test1", repo.Name)
}

func TestServiceGetRepoErrWhenNoOwnerAndNoPublic(t *testing.T) {
	t.Parallel()

	repoID := kallax.NewULID()
	store := &mockStore{repo: &pkg.Repository{ID: repoID, Name: "test1", Visibility: pkg.Private}}
	s := getting.NewService(store)

	userID := "0162eb39-a65e-04a1-7ad9-d663bb49a396"
	u := &pkg.User{ID: userID}
	_, err := s.GetRepo(repoID.String(), u)
	assert.EqualError(t, err, fmt.Sprintf("user %v not authorized to get repo %v", userID, repoID))
	authErr, ok := errors.Cause(err).(authorizationErr)
	assert.True(t, ok)
	assert.True(t, authErr.IsNotAuthorized())
}

func TestServiceGetRepoErrGettingRepoFromStorage(t *testing.T) {
	t.Parallel()

	store := &mockStore{err: errors.New("test")}
	s := getting.NewService(store)

	userID := "0162eb39-a65e-04a1-7ad9-d663bb49a396"
	u := &pkg.User{ID: userID}
	_, err := s.GetRepo("0167c8a5-d308-8692-809d-b1ad4a2d9562", u)
	assert.EqualError(t, err, "could not get repo: test")
}
