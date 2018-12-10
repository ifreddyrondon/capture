package decoder_test

import (
	"testing"
	"time"

	"github.com/ifreddyrondon/capture/pkg"
	"github.com/ifreddyrondon/capture/pkg/repository/decoder"
	"github.com/stretchr/testify/assert"
	"gopkg.in/src-d/go-kallax.v1"
)

func TestRepositoryFromPostRepositoryOK(t *testing.T) {
	t.Parallel()
	name, public, private := "test_repository", "public", "private"

	tt := []struct {
		name     string
		postRepo decoder.PostRepository
		expected pkg.Repository
	}{
		{
			name:     "get repository from postRepository without shared",
			postRepo: decoder.PostRepository{Name: &name, Visibility: nil},
			expected: pkg.Repository{Name: name, Visibility: pkg.Public},
		},
		{
			name:     "get repository from postRepository with shared true",
			postRepo: decoder.PostRepository{Name: &name, Visibility: &public},
			expected: pkg.Repository{Name: name, Visibility: pkg.Public},
		},
		{
			name:     "get repository from postRepository with shared false",
			postRepo: decoder.PostRepository{Name: &name, Visibility: &private},
			expected: pkg.Repository{Name: name, Visibility: pkg.Private},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			repo := tc.postRepo.GetRepository()
			assert.Equal(t, tc.expected.Name, repo.Name)
			assert.Equal(t, tc.expected.Visibility, repo.Visibility)
			// test user fields filled with not default values
			assert.NotEqual(t, kallax.ULID{}, repo.ID)
			assert.NotEqual(t, time.Time{}, repo.CreatedAt)
			assert.NotEqual(t, time.Time{}, repo.UpdatedAt)
		})
	}
}
