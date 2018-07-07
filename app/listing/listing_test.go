package listing_test

import (
	"testing"

	"github.com/ifreddyrondon/capture/app/listing"
	"github.com/ifreddyrondon/capture/app/listing/paging"
	"github.com/ifreddyrondon/capture/app/listing/sorting"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarshalListing(t *testing.T) {
	createdDescSort := sorting.NewSort("created_at_desc", "Created date descending")
	createdAscSort := sorting.NewSort("created_at_asc", "Created date ascendant")

	tt := []struct {
		name     string
		l        listing.Listing
		expected string
	}{
		{
			"given a listing with defaults should marshal only paging",
			listing.Listing{
				Paging: paging.Paging{
					Limit:           paging.DefaultLimit,
					Offset:          paging.DefaultOffset,
					MaxAllowedLimit: paging.DefaultMaxAllowedLimit,
				},
			},
			`{"paging":{"max_allowed_limit":100,"limit":10,"offset":0}}`,
		},
		{
			"given a listing with Paging that includes total should marshal paging with total",
			listing.Listing{
				Paging: paging.Paging{
					Limit:           paging.DefaultLimit,
					Offset:          paging.DefaultOffset,
					MaxAllowedLimit: paging.DefaultMaxAllowedLimit,
					Total:           1000,
				},
			},
			`{"paging":{"max_allowed_limit":100,"limit":10,"offset":0,"total":1000}}`,
		},
		{
			"given a listing with Paging and Sorting should marshal both",
			listing.Listing{
				Paging: paging.Paging{
					Limit:           20,
					Offset:          10,
					MaxAllowedLimit: 50,
				},
				Sorting: sorting.Sorting{
					Sort:      &createdDescSort,
					Available: []sorting.Sort{createdDescSort, createdAscSort},
				},
			},
			`{"paging":{"max_allowed_limit":50,"limit":20,"offset":10},"sort":{"id":"created_at_desc","name":"Created date descending"},"available":[{"id":"created_at_desc","name":"Created date descending"},{"id":"created_at_asc","name":"Created date ascendant"}]}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result, err := tc.l.MarshalJSON()
			require.Nil(t, err)
			assert.Equal(t, tc.expected, string(result))
		})
	}
}
