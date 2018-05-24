package collection_test

import (
	"testing"

	"github.com/ifreddyrondon/gocapture/collection"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalJSON(t *testing.T) {
	t.Parallel()

	expected := collection.Collection{
		Name: "test_collection",
	}

	result := collection.Collection{}
	err := result.UnmarshalJSON([]byte(`{"name":"test_collection"}`))
	require.Nil(t, err)
	assert.Equal(t, expected.Name, result.Name)
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
			[]string{"cannot unmarshal json into valid collection"},
		},
		{
			"empty name",
			[]byte(`{"name":""}`),
			[]string{"name must not be blank"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := collection.Collection{}
			err := result.UnmarshalJSON(tc.payload)
			assert.Error(t, err)
			for _, v := range tc.errs {
				assert.Contains(t, err.Error(), v)
			}
		})
	}
}
