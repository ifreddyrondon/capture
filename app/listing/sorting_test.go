package listing_test

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ifreddyrondon/capture/app/listing"
)

func TestNewSorting(t *testing.T) {
	tt := []struct {
		name       string
		availables []listing.Sort
		result     listing.Sorting
	}{
		{
			"create sorting with only one sorting criteria",
			[]listing.Sort{
				listing.NewSort("created_at_desc", "Created date descending"),
			},
			listing.Sorting{
				Available: []listing.Sort{
					listing.NewSort("created_at_desc", "Created date descending"),
				},
				Sort: listing.NewSort("created_at_desc", "Created date descending"),
			},
		},
		{
			"create sorting with some sorting criteria",
			[]listing.Sort{
				listing.NewSort("created_at_desc", "Created date descending"),
				listing.NewSort("created_at_asc", "Created date ascending"),
				listing.NewSort("updated_at_desc", "Updated date descending"),
				listing.NewSort("updated_at_asc", "Updated date ascending"),
			},
			listing.Sorting{
				Available: []listing.Sort{
					listing.NewSort("created_at_desc", "Created date descending"),
					listing.NewSort("created_at_asc", "Created date ascending"),
					listing.NewSort("updated_at_desc", "Updated date descending"),
					listing.NewSort("updated_at_asc", "Updated date ascending"),
				},
				Sort: listing.NewSort("created_at_desc", "Created date descending"),
			},
		},
		{
			"create sorting with none sorting criteria",
			[]listing.Sort{},
			listing.Sorting{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			s := listing.NewSorting(tc.availables...)
			assert.Equal(t, len(tc.result.Available), len(s.Available))
			assert.Equal(t, tc.result.Sort.ID, s.Sort.ID)
			assert.Equal(t, tc.result.Sort.Name, s.Sort.Name)
		})
	}
}

func TestSortingDecodeOK(t *testing.T) {
	tt := []struct {
		name     string
		params   url.Values
		defaults listing.Sorting
		result   listing.Sorting
	}{
		{
			"given non sort query params present and non sorting criteria",
			map[string][]string{},
			listing.Sorting{},
			listing.Sorting{},
		},
		{
			"given non sort query params present and a sort criteria",
			map[string][]string{},
			listing.NewSorting(listing.NewSort("created_at_desc", "Created date descending")),
			listing.NewSorting(listing.NewSort("created_at_desc", "Created date descending")),
		},
		{
			"given a sort query params present and a sort criteria",
			map[string][]string{"sort": []string{"created_at_desc"}},
			listing.NewSorting(listing.NewSort("created_at_desc", "Created date descending")),
			listing.NewSorting(listing.NewSort("created_at_desc", "Created date descending")),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var s listing.Sorting
			err := s.Decode(tc.params, tc.defaults)
			assert.Nil(t, err)
			assert.Equal(t, tc.result.Sort.ID, s.Sort.ID)
			assert.Equal(t, tc.result.Sort.Name, s.Sort.Name)
			assert.Equal(t, len(tc.result.Available), len(s.Available))
		})
	}
}

func TestSortingDecodeBad(t *testing.T) {
	tt := []struct {
		name     string
		params   url.Values
		defaults listing.Sorting
		err      string
	}{
		{
			"given a sort query not when non sorting criteria",
			map[string][]string{"sort": []string{"foo_desc"}},
			listing.NewSorting(),
			"There's no order criteria with the id foo_desc",
		},
		{
			"given a sort query not when none match sorting criteria",
			map[string][]string{"sort": []string{"foo_desc"}},
			listing.NewSorting(listing.NewSort("created_at_desc", "Created date descending")),
			"There's no order criteria with the id foo_desc",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var s listing.Sorting
			err := s.Decode(tc.params, tc.defaults)
			assert.NotNil(t, err)
			assert.EqualError(t, err, tc.err)
		})
	}
}
