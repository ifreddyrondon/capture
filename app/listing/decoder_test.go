package listing_test

import (
	"net/url"
	"testing"

	"github.com/ifreddyrondon/capture/app/listing"
	"github.com/ifreddyrondon/capture/app/listing/paging"
	"github.com/ifreddyrondon/capture/app/listing/sorting"
	"github.com/stretchr/testify/assert"
)

func TestDecodeOK(t *testing.T) {
	t.Parallel()

	createdDescSort := sorting.NewSort("created_at_desc", "Created date descending")

	tt := []struct {
		name      string
		urlParams url.Values
		opts      []func(*listing.Decoder)
		result    listing.Listing
	}{
		{
			"given none query params and non options should decode paging with defaults",
			map[string][]string{},
			[]func(*listing.Decoder){},
			func() listing.Listing {
				return listing.Listing{
					Paging: paging.Paging{
						Limit:           paging.DefaultLimit,
						Offset:          paging.DefaultOffset,
						MaxAllowedLimit: paging.DefaultMaxAllowedLimit,
					},
					Sorting: sorting.Sorting{},
				}
			}(),
		},
		{
			"given none query params and limit option should decode paging defaults with new limit",
			map[string][]string{},
			[]func(*listing.Decoder){listing.DecodeLimit(50)},
			func() listing.Listing {
				return listing.Listing{
					Paging: paging.Paging{
						Limit:           50,
						Offset:          paging.DefaultOffset,
						MaxAllowedLimit: paging.DefaultMaxAllowedLimit,
					},
					Sorting: sorting.Sorting{},
				}
			}(),
		},
		{
			"given offset and limit when limit > maxAllowed default and maxAllowed option should decode paging with offset and limit upper the default",
			map[string][]string{"offset": []string{"1"}, "limit": []string{"105"}},
			[]func(*listing.Decoder){listing.DecodeMaxAllowedLimit(110)},
			func() listing.Listing {
				return listing.Listing{
					Paging: paging.Paging{
						Limit:           105,
						Offset:          1,
						MaxAllowedLimit: 110,
					},
					Sorting: sorting.Sorting{},
				}
			}(),
		},
		{
			"given a sort params and sort criteria option should decode sorting with sort and availables criteria",
			map[string][]string{"sort": []string{"created_at_desc"}},
			[]func(*listing.Decoder){listing.DecodeSortCriterias(sorting.NewSort("created_at_desc", "Created date descending"))},
			func() listing.Listing {
				return listing.Listing{
					Paging: paging.Paging{
						Limit:           paging.DefaultLimit,
						Offset:          paging.DefaultOffset,
						MaxAllowedLimit: paging.DefaultMaxAllowedLimit,
					},
					Sorting: sorting.Sorting{
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
		err       string
	}{
		{
			"given a not number limit param should return an error when decode paging",
			map[string][]string{"limit": []string{"a"}},
			"invalid limit value, must be a number",
		},
		{
			"given a sort query when non sorting criteria",
			map[string][]string{"sort": []string{"a"}},
			"there's no order criteria with the id a",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var l listing.Listing
			err := listing.NewDecoder(tc.urlParams).Decode(&l)
			assert.NotNil(t, err)
			assert.EqualError(t, err, tc.err)
		})
	}
}
