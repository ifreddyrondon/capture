package creating_test

import (
	"testing"

	"github.com/ifreddyrondon/capture/pkg/creating"
	"github.com/ifreddyrondon/capture/pkg/domain"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"gopkg.in/src-d/go-kallax.v1"
)

type mockStore struct {
	err error
}

func (m *mockStore) SaveRepo(repo *domain.Repository) error { return m.err }

func string2pointer(v string) *string { return &v }

func TestServiceCreateRepoOK(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name     string
		payl     creating.Payload
		expected creating.Repository
	}{
		{
			name:     "given a only name payload should return a repo with visibility public",
			payl:     creating.Payload{Name: string2pointer("test")},
			expected: creating.Repository{Name: "test", Visibility: "public"},
		},
		{
			name:     "given a payload with name and visibility should return a repo with visibility private",
			payl:     creating.Payload{Name: string2pointer("test"), Visibility: string2pointer("private")},
			expected: creating.Repository{Name: "test", Visibility: "private"},
		},
	}

	owner := &domain.User{ID: kallax.NewULID()}
	s := creating.NewService(&mockStore{})

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			repo, err := s.CreateRepo(owner, tc.payl)
			assert.Nil(t, err)

			assert.NotNil(t, repo.ID)
			assert.Equal(t, tc.expected.Name, repo.Name)
			assert.Equal(t, tc.expected.Visibility, repo.Visibility)
			assert.NotNil(t, repo.CreatedAt)
			assert.NotNil(t, repo.UpdatedAt)
		})
	}
}

func TestServiceCreateRepoErrWhenSaving(t *testing.T) {
	t.Parallel()
	s := creating.NewService(&mockStore{err: errors.New("test")})

	owner := &domain.User{ID: kallax.NewULID()}
	payl := creating.Payload{Name: string2pointer("test")}
	_, err := s.CreateRepo(owner, payl)
	assert.EqualError(t, err, "could not save repo: test")
}
