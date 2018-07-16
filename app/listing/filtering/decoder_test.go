package filtering_test

import (
	"net/url"
	"testing"

	"github.com/ifreddyrondon/capture/app/listing/filtering"
	"github.com/stretchr/testify/assert"
)

func TestDecodeOK(t *testing.T) {
	t.Parallel()

	vNew := filtering.NewValue("new", "New")
	vUsed := filtering.NewValue("used", "Used")
	text := filtering.NewText("condition", "test", vNew, vUsed)
	boolean := filtering.NewBoolean("shared", "test", "shared", "private")

	tt := []struct {
		name      string
		urlParams url.Values
		decoders  []filtering.FilterDecoder
		result    filtering.Filtering
	}{
		{
			"given none query params and non decoders should decode empty Filtering",
			map[string][]string{},
			[]filtering.FilterDecoder{},
			filtering.Filtering{},
		},
		{
			"given non filter query params present and one decoder should decode empty filter with all availables",
			map[string][]string{},
			[]filtering.FilterDecoder{text},
			filtering.Filtering{
				Filters: []filtering.Filter{},
				Available: []filtering.Filter{
					filtering.Filter{
						ID:     "condition",
						Name:   "test",
						Type:   "text",
						Values: []filtering.Value{vNew, vUsed},
					},
				},
			},
		},
		{
			"given non filter query params present and some decoders should decode empty filter with all filters availables",
			map[string][]string{},
			[]filtering.FilterDecoder{text, boolean},
			filtering.Filtering{
				Filters: []filtering.Filter{},
				Available: []filtering.Filter{
					filtering.Filter{
						ID:     "condition",
						Name:   "test",
						Type:   "text",
						Values: []filtering.Value{vNew, vUsed},
					},
					filtering.Filter{
						ID:   "shared",
						Name: "test",
						Type: "boolean",
						Values: []filtering.Value{
							filtering.NewValue("true", "shared"),
							filtering.NewValue("false", "private"),
						},
					},
				},
			},
		},
		{
			"given condition query params with one value and one decoder should decode condition filter with the value and with all availables",
			map[string][]string{"condition": []string{"new"}},
			[]filtering.FilterDecoder{text},
			filtering.Filtering{
				Filters: []filtering.Filter{
					filtering.Filter{
						ID:     "condition",
						Name:   "test",
						Type:   "text",
						Values: []filtering.Value{vNew},
					},
				},
				Available: []filtering.Filter{
					filtering.Filter{
						ID:     "condition",
						Name:   "test",
						Type:   "text",
						Values: []filtering.Value{vNew, vUsed},
					},
				},
			},
		},
		{
			"given condition query params with two value and one decoder should decode condition filter with the first value and with filter availables",
			map[string][]string{"condition": []string{"new"}},
			[]filtering.FilterDecoder{text},
			filtering.Filtering{
				Filters: []filtering.Filter{
					filtering.Filter{
						ID:     "condition",
						Name:   "test",
						Type:   "text",
						Values: []filtering.Value{vNew},
					},
				},
				Available: []filtering.Filter{
					filtering.Filter{
						ID:     "condition",
						Name:   "test",
						Type:   "text",
						Values: []filtering.Value{vNew, vUsed},
					},
				},
			},
		},
		{
			"given condition query params with one value and some decoders should decode condition filter with the value and with all filters availables",
			map[string][]string{"condition": []string{"new"}},
			[]filtering.FilterDecoder{text, boolean},
			filtering.Filtering{
				Filters: []filtering.Filter{
					filtering.Filter{
						ID:     "condition",
						Name:   "test",
						Type:   "text",
						Values: []filtering.Value{vNew},
					},
				},
				Available: []filtering.Filter{
					filtering.Filter{
						ID:     "condition",
						Name:   "test",
						Type:   "text",
						Values: []filtering.Value{vNew, vUsed},
					},
					filtering.Filter{
						ID:   "shared",
						Name: "test",
						Type: "boolean",
						Values: []filtering.Value{
							filtering.NewValue("true", "shared"),
							filtering.NewValue("false", "private"),
						},
					},
				},
			},
		},
		{
			"given condition query params with two values and some decoders should decode all filters selected with the value and with all filters availables",
			map[string][]string{"condition": []string{"new"}, "shared": []string{"true"}},
			[]filtering.FilterDecoder{text, boolean},
			filtering.Filtering{
				Filters: []filtering.Filter{
					filtering.Filter{
						ID:     "condition",
						Name:   "test",
						Type:   "text",
						Values: []filtering.Value{vNew},
					},
					filtering.Filter{
						ID:     "shared",
						Name:   "test",
						Type:   "boolean",
						Values: []filtering.Value{filtering.NewValue("true", "shared")},
					},
				},
				Available: []filtering.Filter{
					filtering.Filter{
						ID:     "condition",
						Name:   "test",
						Type:   "text",
						Values: []filtering.Value{vNew, vUsed},
					},
					filtering.Filter{
						ID:   "shared",
						Name: "test",
						Type: "boolean",
						Values: []filtering.Value{
							filtering.NewValue("true", "shared"),
							filtering.NewValue("false", "private"),
						},
					},
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var f filtering.Filtering
			err := filtering.NewDecoder(tc.urlParams, tc.decoders...).Decode(&f)
			assert.Nil(t, err)
			for i, v := range f.Filters {
				assert.Equal(t, tc.result.Filters[i], v)
			}
			for i, v := range f.Available {
				assert.Equal(t, tc.result.Available[i], v)
			}
		})
	}
}

func TestDecodeMissing(t *testing.T) {
	t.Parallel()

	vNew := filtering.NewValue("new", "New")
	vUsed := filtering.NewValue("used", "Used")
	text := filtering.NewText("condition", "test", vNew, vUsed)
	boolean := filtering.NewBoolean("shared", "test", "shared", "private")

	tt := []struct {
		name      string
		urlParams url.Values
		decoders  []filtering.FilterDecoder
		result    filtering.Filtering
	}{
		{
			"given condition query params with one missing value and one decoder should decode empty filter and with all availables",
			map[string][]string{"condition": []string{"some"}},
			[]filtering.FilterDecoder{text},
			filtering.Filtering{
				Filters: []filtering.Filter{},
				Available: []filtering.Filter{
					filtering.Filter{
						ID:     "condition",
						Name:   "test",
						Type:   "text",
						Values: []filtering.Value{vNew, vUsed},
					},
				},
			},
		},
		{
			"given query params with no match and one decoder should decode empty filter and with all availables",
			map[string][]string{"condition": []string{"some"}},
			[]filtering.FilterDecoder{text},
			filtering.Filtering{
				Filters: []filtering.Filter{},
				Available: []filtering.Filter{
					filtering.Filter{
						ID:     "condition",
						Name:   "test",
						Type:   "text",
						Values: []filtering.Value{vNew, vUsed},
					},
				},
			},
		},
		{
			"given somes query params with no match and somes decoder should decode empty filter and with all availables",
			map[string][]string{"foo": []string{"new"}, "faa": []string{"true"}},
			[]filtering.FilterDecoder{text, boolean},
			filtering.Filtering{
				Filters: []filtering.Filter{},
				Available: []filtering.Filter{
					filtering.Filter{
						ID:     "condition",
						Name:   "test",
						Type:   "text",
						Values: []filtering.Value{vNew, vUsed},
					},
					filtering.Filter{
						ID:   "shared",
						Name: "test",
						Type: "boolean",
						Values: []filtering.Value{
							filtering.NewValue("true", "shared"),
							filtering.NewValue("false", "private"),
						},
					},
				},
			},
		},
		{
			"given somes query params with filter but no match values and somes decoder should decode empty filter and with all availables",
			map[string][]string{"condition": []string{"foo"}, "shared": []string{"faa"}},
			[]filtering.FilterDecoder{text, boolean},
			filtering.Filtering{
				Filters: []filtering.Filter{},
				Available: []filtering.Filter{
					filtering.Filter{
						ID:     "condition",
						Name:   "test",
						Type:   "text",
						Values: []filtering.Value{vNew, vUsed},
					},
					filtering.Filter{
						ID:   "shared",
						Name: "test",
						Type: "boolean",
						Values: []filtering.Value{
							filtering.NewValue("true", "shared"),
							filtering.NewValue("false", "private"),
						},
					},
				},
			},
		},
		{
			"given somes query params with filter and only one with match values and somes decoder should decode filter with the match and with all availables",
			map[string][]string{"condition": []string{"foo"}, "shared": []string{"false"}},
			[]filtering.FilterDecoder{text, boolean},
			filtering.Filtering{
				Filters: []filtering.Filter{
					filtering.Filter{
						ID:     "shared",
						Name:   "test",
						Type:   "boolean",
						Values: []filtering.Value{filtering.NewValue("false", "private")},
					},
				},
				Available: []filtering.Filter{
					filtering.Filter{
						ID:     "condition",
						Name:   "test",
						Type:   "text",
						Values: []filtering.Value{vNew, vUsed},
					},
					filtering.Filter{
						ID:   "shared",
						Name: "test",
						Type: "boolean",
						Values: []filtering.Value{
							filtering.NewValue("true", "shared"),
							filtering.NewValue("false", "private"),
						},
					},
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var f filtering.Filtering
			err := filtering.NewDecoder(tc.urlParams, tc.decoders...).Decode(&f)
			assert.Nil(t, err)
			for i, v := range f.Filters {
				assert.Equal(t, tc.result.Filters[i], v)
			}
			for i, v := range f.Available {
				assert.Equal(t, tc.result.Available[i], v)
			}
		})
	}
}
