package listing_test

import (
	"net/url"
	"testing"

	"github.com/ifreddyrondon/capture/app/listing"
	"github.com/ifreddyrondon/capture/app/listing/filtering"
	"github.com/ifreddyrondon/capture/app/listing/paging"
	"github.com/ifreddyrondon/capture/app/listing/sorting"
	"github.com/stretchr/testify/assert"
)

func TestDecodeOK(t *testing.T) {
	t.Parallel()

	createdDescSort := sorting.NewSort("created_at_desc", "Created date descending")
	vNew := filtering.NewValue("new", "New")
	vUsed := filtering.NewValue("used", "Used")
	vTrue := filtering.NewValue("true", "shared")
	vFalse := filtering.NewValue("false", "private")
	text := filtering.NewText("condition", "test", vNew, vUsed)
	boolean := filtering.NewBoolean("shared", "test", "shared", "private")

	tt := []struct {
		name      string
		urlParams url.Values
		opts      []func(*listing.Decoder)
		result    listing.Listing
	}{
		{
			"given none query params with non options should decode paging with defaults",
			map[string][]string{},
			[]func(*listing.Decoder){},
			func() listing.Listing {
				return listing.Listing{
					Paging: paging.Paging{
						Limit:           paging.DefaultLimit,
						Offset:          paging.DefaultOffset,
						MaxAllowedLimit: paging.DefaultMaxAllowedLimit,
					},
				}
			}(),
		},
		{
			"given none query params with limit option should decode paging defaults with new limit",
			map[string][]string{},
			[]func(*listing.Decoder){listing.DecodeLimit(50)},
			func() listing.Listing {
				return listing.Listing{
					Paging: paging.Paging{
						Limit:           50,
						Offset:          paging.DefaultOffset,
						MaxAllowedLimit: paging.DefaultMaxAllowedLimit,
					},
				}
			}(),
		},
		{
			"given offset and limit when limit > maxAllowed default with maxAllowed option should decode paging with offset and limit upper the default",
			map[string][]string{"offset": []string{"1"}, "limit": []string{"105"}},
			[]func(*listing.Decoder){listing.DecodeMaxAllowedLimit(110)},
			func() listing.Listing {
				return listing.Listing{
					Paging: paging.Paging{
						Limit:           105,
						Offset:          1,
						MaxAllowedLimit: 110,
					},
				}
			}(),
		},
		{
			"given a sort params with sort criteria and filter criteria should decode sorting with availables criteria also decode filtering with only availables",
			map[string][]string{"sort": []string{"created_at_desc"}},
			[]func(*listing.Decoder){
				listing.DecodeSort(createdDescSort),
				listing.DecodeFilter(text),
			},
			func() listing.Listing {
				return listing.Listing{
					Paging: paging.Paging{
						Limit:           paging.DefaultLimit,
						Offset:          paging.DefaultOffset,
						MaxAllowedLimit: paging.DefaultMaxAllowedLimit,
					},
					Sorting: &sorting.Sorting{
						Sort:      &createdDescSort,
						Available: []sorting.Sort{createdDescSort},
					},
					Filtering: &filtering.Filtering{
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
				}
			}(),
		},
		{
			"given a sort and filter params with sort criteria and filter criteria should decode sorting with availables criteria also decode filtering with filter and availables",
			map[string][]string{"sort": []string{"created_at_desc"}, "condition": []string{"new"}},
			[]func(*listing.Decoder){
				listing.DecodeSort(createdDescSort),
				listing.DecodeFilter(text, boolean),
			},
			func() listing.Listing {
				return listing.Listing{
					Paging: paging.Paging{
						Limit:           paging.DefaultLimit,
						Offset:          paging.DefaultOffset,
						MaxAllowedLimit: paging.DefaultMaxAllowedLimit,
					},
					Sorting: &sorting.Sorting{
						Sort:      &createdDescSort,
						Available: []sorting.Sort{createdDescSort},
					},
					Filtering: &filtering.Filtering{
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
								ID:     "shared",
								Name:   "test",
								Type:   "boolean",
								Values: []filtering.Value{vTrue, vFalse},
							},
						},
					},
				}
			}(),
		},
		{
			"given none params with one filter criteria should decode filtering with empty filters and availables criteria",
			map[string][]string{"sort": []string{"created_at_desc"}},
			[]func(*listing.Decoder){listing.DecodeSort(createdDescSort)},
			func() listing.Listing {
				return listing.Listing{
					Paging: paging.Paging{
						Limit:           paging.DefaultLimit,
						Offset:          paging.DefaultOffset,
						MaxAllowedLimit: paging.DefaultMaxAllowedLimit,
					},
					Sorting: &sorting.Sorting{
						Sort:      &createdDescSort,
						Available: []sorting.Sort{createdDescSort},
					},
				}
			}(),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var l listing.Listing
			err := listing.NewDecoder(tc.urlParams, tc.opts...).Decode(&l)
			assert.Nil(t, err)
			assert.Equal(t, l.Paging, tc.result.Paging)
			assert.Equal(t, l.Sorting, tc.result.Sorting)
		})
	}
}

func TestDecodeFails(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name      string
		urlParams url.Values
		opts      []func(*listing.Decoder)
		err       string
	}{
		{
			"given a not number limit param should return an error when decode paging",
			map[string][]string{"limit": []string{"a"}},
			[]func(*listing.Decoder){},
			"invalid limit value, must be a number",
		},
		{
			"given a sort query when non match sorting criteria",
			map[string][]string{"sort": []string{"a"}},
			[]func(*listing.Decoder){listing.DecodeSort(sorting.NewSort("created_at_desc", "Created date descending"))},
			"there's no order criteria with the id a",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var l listing.Listing
			err := listing.NewDecoder(tc.urlParams, tc.opts...).Decode(&l)
			assert.NotNil(t, err)
			assert.EqualError(t, err, tc.err)
		})
	}
}
