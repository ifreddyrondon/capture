package filtering_test

import (
	"net/url"
	"testing"

	"github.com/ifreddyrondon/capture/app/listing/filtering"
	"github.com/stretchr/testify/assert"
)

func TestTextPresentOK(t *testing.T) {
	t.Parallel()

	vNew := filtering.NewValue("new", "New")
	vUsed := filtering.NewValue("used", "Used")

	tt := []struct {
		name     string
		decoder  *filtering.Text
		params   url.Values
		expected *filtering.Filter
	}{
		{
			"should return new value when param with condition with new value",
			filtering.NewText("condition", "test", vNew),
			map[string][]string{"condition": []string{"new"}},
			&filtering.Filter{
				ID:     "condition",
				Name:   "test",
				Type:   "text",
				Values: []filtering.Value{filtering.NewValue("new", "New")},
			},
		},
		{
			"should return used value when param with condition with used value",
			filtering.NewText("condition", "test", vUsed),
			map[string][]string{"condition": []string{"used"}},
			&filtering.Filter{
				ID:     "condition",
				Name:   "test",
				Type:   "text",
				Values: []filtering.Value{filtering.NewValue("used", "Used")},
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

func TestTextPresentFails(t *testing.T) {
	t.Parallel()

	vNew := filtering.NewValue("new", "New")
	tt := []struct {
		name    string
		decoder *filtering.Text
		params  url.Values
	}{
		{
			"should return nil when not value found",
			filtering.NewText("condition", "test", vNew),
			map[string][]string{"condition": []string{"abc"}},
		},
		{
			"should return nil when not params found",
			filtering.NewText("condition", "test", vNew),
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

func TestTextWithValues(t *testing.T) {
	t.Parallel()

	values := []filtering.Value{
		filtering.NewValue("new", "New"),
		filtering.NewValue("used", "Used"),
	}
	decoder := filtering.NewText("condition", "test", values...)
	expected := &filtering.Filter{
		ID:     "condition",
		Name:   "test",
		Type:   "text",
		Values: values,
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
