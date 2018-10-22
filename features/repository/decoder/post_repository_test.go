package decoder_test

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/ifreddyrondon/capture/features"
	"github.com/ifreddyrondon/capture/features/repository/decoder"
	"github.com/stretchr/testify/assert"
	"gopkg.in/src-d/go-kallax.v1"
)

func TestDecodePostRepositoryOK(t *testing.T) {
	t.Parallel()

	name, sharedTrue, sharedFalse := "test_repository", true, false

	tt := []struct {
		name     string
		body     string
		expected decoder.PostRepository
	}{
		{
			name:     "decode repo without shared",
			body:     `{"name":"test_repository"}`,
			expected: decoder.PostRepository{Name: &name, Shared: nil},
		},
		{
			name:     "decode repo with shared true",
			body:     `{"name":"test_repository","shared":true}`,
			expected: decoder.PostRepository{Name: &name, Shared: &sharedTrue},
		},
		{
			name:     "decode repo with shared true",
			body:     `{"name":"test_repository","shared":false}`,
			expected: decoder.PostRepository{Name: &name, Shared: &sharedFalse},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			r, _ := http.NewRequest("POST", "/", strings.NewReader(tc.body))

			var repo decoder.PostRepository
			err := decoder.Decode(r, &repo)
			assert.Nil(t, err)
			assert.Equal(t, tc.expected.Name, repo.Name)
			assert.Equal(t, tc.expected.Shared, repo.Shared)
		})
	}
}

func TestDecodePostRepositoryError(t *testing.T) {
	t.Parallel()
	tt := []struct {
		name string
		body string
		err  string
	}{
		{
			name: "decode repository name missing",
			body: `{}`,
			err:  "name must not be blank",
		},
		{
			name: "decode repository name empty",
			body: `{"name":""}`,
			err:  "name must not be blank",
		},
		{
			name: "decode repository name empty v2",
			body: `{"name":"   "}`,
			err:  "name must not be blank",
		},
		{
			name: "invalid repository payload",
			body: `.`,
			err:  "cannot unmarshal json into valid repository",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			r, _ := http.NewRequest("POST", "/", strings.NewReader(tc.body))

			var repo decoder.PostRepository
			err := decoder.Decode(r, &repo)
			assert.EqualError(t, err, tc.err)
		})
	}
}

func TestRepositoryFromPostRepositoryOK(t *testing.T) {
	name, sharedTrue, sharedFalse := "test_repository", true, false
	t.Parallel()
	tt := []struct {
		name     string
		postRepo decoder.PostRepository
		expected features.Repository
	}{
		{
			name:     "get repository from postRepository without shared",
			postRepo: decoder.PostRepository{Name: &name, Shared: nil},
			expected: features.Repository{Name: name, Shared: true},
		},
		{
			name:     "get repository from postRepository with shared true",
			postRepo: decoder.PostRepository{Name: &name, Shared: &sharedTrue},
			expected: features.Repository{Name: name, Shared: true},
		},
		{
			name:     "get repository from postRepository with shared false",
			postRepo: decoder.PostRepository{Name: &name, Shared: &sharedFalse},
			expected: features.Repository{Name: name, Shared: false},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			repo := tc.postRepo.GetRepository()
			assert.Equal(t, tc.expected.Name, repo.Name)
			assert.Equal(t, tc.expected.Shared, repo.Shared)
			// test user fields filled with not default values
			assert.NotEqual(t, kallax.ULID{}, repo.ID)
			assert.NotEqual(t, time.Time{}, repo.CreatedAt)
			assert.NotEqual(t, time.Time{}, repo.UpdatedAt)
		})
	}
}
