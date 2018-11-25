package decoder_test

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/ifreddyrondon/capture/pkg"
	"github.com/ifreddyrondon/capture/pkg/repository/decoder"
	"github.com/stretchr/testify/assert"
	"gopkg.in/src-d/go-kallax.v1"
)

func TestDecodePostRepositoryOK(t *testing.T) {
	t.Parallel()
	name, public, private := "test_repository", "public", "private"

	tt := []struct {
		name     string
		body     string
		expected decoder.PostRepository
	}{
		{
			name:     "decode repo without shared",
			body:     `{"name":"test_repository"}`,
			expected: decoder.PostRepository{Name: &name, Visibility: nil},
		},
		{
			name:     "decode repo with shared true",
			body:     `{"name":"test_repository","visibility":"public"}`,
			expected: decoder.PostRepository{Name: &name, Visibility: &public},
		},
		{
			name:     "decode repo with shared true",
			body:     `{"name":"test_repository","visibility":"private"}`,
			expected: decoder.PostRepository{Name: &name, Visibility: &private},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			r, _ := http.NewRequest("POST", "/", strings.NewReader(tc.body))

			var repo decoder.PostRepository
			err := decoder.Decode(r, &repo)
			assert.Nil(t, err)
			assert.Equal(t, tc.expected.Name, repo.Name)
			assert.Equal(t, tc.expected.Visibility, repo.Visibility)
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
			name: "decode repository with not allowed visibility",
			body: `{"name":"foo","visibility":"protected"}`,
			err:  "not allowed visibility type. it Could be one of public, or private. Default: public",
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
