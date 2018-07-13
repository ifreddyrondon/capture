package filtering_test

import (
	"net/url"
	"testing"

	"github.com/ifreddyrondon/capture/app/listing/filtering"
	"github.com/stretchr/testify/assert"
)

func TestNewBooleanBuilder(t *testing.T) {
	f := filtering.NewBooleanBuilder("shared", "shared collections filter", "shared collection", "private collection")
	assert.Equal(t, "shared", f.ID)
	assert.Equal(t, "shared collections filter", f.Name)
}

func TestBooleanBuilderValidateOK(t *testing.T) {
	tt := []struct {
		name     string
		builder  *filtering.BooleanBuilder
		params   url.Values
		expected *filtering.Filter
	}{
		{
			"should return true value when param with true value",
			filtering.NewBooleanBuilder("shared", "test", "shared", "private"),
			map[string][]string{"shared": []string{"true"}},
			&filtering.Filter{
				FilterID: filtering.NewFilterID("shared", "test"),
				Type:     "boolean",
				Values:   []filtering.FilterValue{filtering.NewFilterValue("true", "shared")},
			},
		},
		{
			"should return false value when param with false value",
			filtering.NewBooleanBuilder("shared", "test", "shared", "private"),
			map[string][]string{"shared": []string{"false"}},
			&filtering.Filter{
				FilterID: filtering.NewFilterID("shared", "test"),
				Type:     "boolean",
				Values:   []filtering.FilterValue{filtering.NewFilterValue("false", "private")},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.builder.Validate(tc.params)
			assert.Equal(t, tc.expected.ID, result.ID)
			assert.Equal(t, tc.expected.Name, result.Name)
			assert.Equal(t, tc.expected.Type, result.Type)
			assert.Equal(t, len(tc.expected.Values), len(result.Values))
			assert.Equal(t, tc.expected.Values[0].ID, result.Values[0].ID)
			assert.Equal(t, tc.expected.Values[0].Name, result.Values[0].Name)
		})
	}
}

func TestBooleanBuilderValidateFails(t *testing.T) {
	tt := []struct {
		name    string
		builder *filtering.BooleanBuilder
		params  url.Values
	}{
		{
			"should return nil when not value found",
			filtering.NewBooleanBuilder("shared", "test", "shared", "private"),
			map[string][]string{"shared": []string{"abc"}},
		},
		{
			"should return nil when not params found",
			filtering.NewBooleanBuilder("shared", "test", "shared", "private"),
			map[string][]string{"foo": []string{"abc"}},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.builder.Validate(tc.params)
			assert.Nil(t, result)
		})
	}
}

func TestBooleanBuilderWithValues(t *testing.T) {
	builder := filtering.NewBooleanBuilder("shared", "test", "shared", "private")
	expected := &filtering.Filter{
		FilterID: filtering.NewFilterID("shared", "test"),
		Type:     "boolean",
		Values: []filtering.FilterValue{
			filtering.NewFilterValue("true", "shared"),
			filtering.NewFilterValue("false", "private"),
		},
	}
	result := builder.WithValues()
	assert.Equal(t, expected.ID, result.ID)
	assert.Equal(t, expected.Name, result.Name)
	assert.Equal(t, expected.Type, result.Type)
	assert.Equal(t, len(expected.Values), len(result.Values))
	assert.Equal(t, expected.Values[0].ID, result.Values[0].ID)
	assert.Equal(t, expected.Values[0].Name, result.Values[0].Name)
	assert.Equal(t, expected.Values[1].ID, result.Values[1].ID)
	assert.Equal(t, expected.Values[1].Name, result.Values[1].Name)
}
