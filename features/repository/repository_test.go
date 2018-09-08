package repository_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/ifreddyrondon/capture/features/repository"

	kallax "gopkg.in/src-d/go-kallax.v1"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalJSON(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name     string
		payload  []byte
		expected repository.Repository
	}{
		{
			"just name",
			[]byte(`{"name":"test_repository"}`),
			repository.Repository{Name: "test_repository", Shared: false},
		},
		{
			"name and shared",
			[]byte(`{"name":"test_repository", "shared": true}`),
			repository.Repository{Name: "test_repository", Shared: true},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := repository.Repository{}
			err := result.UnmarshalJSON(tc.payload)
			require.Nil(t, err)
			assert.Equal(t, tc.expected.Name, result.Name)
			assert.Equal(t, tc.expected.Shared, result.Shared)
		})
	}
}

func TestUnmarshalJSONFail(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name    string
		payload []byte
		errs    []string
	}{
		{
			"invalid payload",
			[]byte(`{`),
			[]string{"cannot unmarshal json into valid repository"},
		},
		{
			"empty name",
			[]byte(`{"name":""}`),
			[]string{"name must not be blank"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := repository.Repository{}
			err := result.UnmarshalJSON(tc.payload)
			assert.Error(t, err)
			for _, v := range tc.errs {
				assert.Contains(t, err.Error(), v)
			}
		})
	}
}

func TestMarshalJSON(t *testing.T) {
	t.Parallel()

	d, _ := time.Parse(time.RFC3339, "1989-12-26T06:01:00.00Z")

	expected := []byte(`{"id":"0162eb39-a65e-04a1-7ad9-d663bb49a396","name":"test","current_branch":"","shared":true,"createdAt":"1989-12-26T06:01:00Z","updatedAt":"1989-12-26T06:01:00Z","owner":"00000000-0000-0000-0000-000000000000"}`)
	c := repository.Repository{
		Name: "test",
		ID: func() kallax.ULID {
			id, _ := kallax.NewULIDFromText("0162eb39-a65e-04a1-7ad9-d663bb49a396")
			return id
		}(),
		Shared:    true,
		CreatedAt: d,
		UpdatedAt: d,
	}

	result, err := json.Marshal(c)
	require.Nil(t, err)
	assert.Equal(t, expected, result)
}
