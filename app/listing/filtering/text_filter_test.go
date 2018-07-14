package filtering_test

import (
	"net/url"
	"testing"

	"github.com/ifreddyrondon/capture/app/listing/filtering"
	"github.com/stretchr/testify/assert"
)

func TestNewTextDecoder(t *testing.T) {
	vNew := filtering.NewFilterValue("new", "New")
	vUsed := filtering.NewFilterValue("used", "Used")
	f := filtering.NewTextDecoder("condition", "Condición", vNew, vUsed)
	assert.Equal(t, "condition", f.ID)
	assert.Equal(t, "Condición", f.Name)
}

func TestTextDecoderPresentOK(t *testing.T) {
	vNew := filtering.NewFilterValue("new", "New")
	vUsed := filtering.NewFilterValue("used", "Used")

	tt := []struct {
		name     string
		decoder  *filtering.TextDecoder
		params   url.Values
		expected *filtering.Filter
	}{
		{
			"should return new value when param with condition with new value",
			filtering.NewTextDecoder("condition", "test", vNew),
			map[string][]string{"condition": []string{"new"}},
			&filtering.Filter{
				FilterID: filtering.NewFilterID("condition", "test"),
				Type:     "text",
				Values:   []filtering.FilterValue{filtering.NewFilterValue("new", "New")},
			},
		},
		{
			"should return used value when param with condition with used value",
			filtering.NewTextDecoder("condition", "test", vUsed),
			map[string][]string{"condition": []string{"used"}},
			&filtering.Filter{
				FilterID: filtering.NewFilterID("condition", "test"),
				Type:     "text",
				Values:   []filtering.FilterValue{filtering.NewFilterValue("used", "Used")},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.decoder.Present(tc.params)
			assert.Equal(t, tc.expected.ID, result.ID)
			assert.Equal(t, tc.expected.Name, result.Name)
			assert.Equal(t, tc.expected.Type, result.Type)
			assert.Equal(t, len(tc.expected.Values), len(result.Values))
			assert.Equal(t, tc.expected.Values[0].ID, result.Values[0].ID)
			assert.Equal(t, tc.expected.Values[0].Name, result.Values[0].Name)
		})
	}
}

func TestTextDecoderPresentFails(t *testing.T) {
	vNew := filtering.NewFilterValue("new", "New")
	tt := []struct {
		name    string
		decoder *filtering.TextDecoder
		params  url.Values
	}{
		{
			"should return nil when not value found",
			filtering.NewTextDecoder("condition", "test", vNew),
			map[string][]string{"condition": []string{"abc"}},
		},
		{
			"should return nil when not params found",
			filtering.NewTextDecoder("condition", "test", vNew),
			map[string][]string{"foo": []string{"abc"}},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.decoder.Present(tc.params)
			assert.Nil(t, result)
		})
	}
}

func TestTextDecoderWithValues(t *testing.T) {
	values := []filtering.FilterValue{
		filtering.NewFilterValue("new", "New"),
		filtering.NewFilterValue("used", "Used"),
	}
	decoder := filtering.NewTextDecoder("condition", "test", values...)
	expected := &filtering.Filter{
		FilterID: filtering.NewFilterID("condition", "test"),
		Type:     "text",
		Values:   values,
	}
	result := decoder.WithValues()
	assert.Equal(t, expected.ID, result.ID)
	assert.Equal(t, expected.Name, result.Name)
	assert.Equal(t, expected.Type, result.Type)
	assert.Equal(t, len(expected.Values), len(result.Values))
	for i, v := range result.Values {
		assert.Equal(t, expected.Values[i].ID, v.ID)
		assert.Equal(t, expected.Values[i].Name, v.Name)
	}
}
