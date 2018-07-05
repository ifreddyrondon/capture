package sorting_test

import (
	"net/url"
	"testing"

	"github.com/ifreddyrondon/capture/app/listing/sorting"
	"github.com/stretchr/testify/assert"
)

func TestDecodeOK(t *testing.T) {
	tt := []struct {
		name      string
		urlParams url.Values
		criterias []sorting.Sort
		result    sorting.Sorting
	}{
		{
			"given none query params and non criterias should decode empty Sorting",
			map[string][]string{},
			[]sorting.Sort{},
			sorting.Sorting{},
		},
		{
			"given non sort query params present and one sort criteria",
			map[string][]string{},
			[]sorting.Sort{
				sorting.NewSort("created_at_desc", "Created date descending"),
			},
			sorting.Sorting{
				Sort: sorting.NewSort("created_at_desc", "Created date descending"),
				Available: []sorting.Sort{
					sorting.NewSort("created_at_desc", "Created date descending"),
				},
			},
		},
		{
			"given non sort query params present and one some sort criteria",
			map[string][]string{},
			[]sorting.Sort{
				sorting.NewSort("created_at_desc", "Created date descending"),
				sorting.NewSort("created_at_asc", "Created date ascendant"),
			},
			sorting.Sorting{
				Sort: sorting.NewSort("created_at_desc", "Created date descending"),
				Available: []sorting.Sort{
					sorting.NewSort("created_at_desc", "Created date descending"),
					sorting.NewSort("created_at_asc", "Created date ascendant"),
				},
			},
		},
		{
			"given created_at_desc sort query params present and one some sort criteria",
			map[string][]string{"sort": []string{"created_at_asc"}},
			[]sorting.Sort{
				sorting.NewSort("created_at_desc", "Created date descending"),
				sorting.NewSort("created_at_asc", "Created date ascendant"),
			},
			sorting.Sorting{
				Sort: sorting.NewSort("created_at_asc", "Created date ascendant"),
				Available: []sorting.Sort{
					sorting.NewSort("created_at_desc", "Created date descending"),
					sorting.NewSort("created_at_asc", "Created date ascendant"),
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var s sorting.Sorting
			err := sorting.NewDecoder(tc.urlParams, tc.criterias...).Decode(&s)
			assert.Nil(t, err)
			assert.Equal(t, tc.result.Sort.ID, s.Sort.ID)
			assert.Equal(t, tc.result.Sort.Name, s.Sort.Name)
			assert.Equal(t, len(tc.result.Available), len(s.Available))
		})
	}
}

func TestSortingDecodeBad(t *testing.T) {
	tt := []struct {
		name      string
		urlParams url.Values
		criterias []sorting.Sort
		err       string
	}{
		{
			"given a sort query when non sorting criteria",
			map[string][]string{"sort": []string{"foo_desc"}},
			[]sorting.Sort{},
			"there's no order criteria with the id foo_desc",
		},
		{
			"given a sort query when none match sorting criteria",
			map[string][]string{"sort": []string{"foo_desc"}},
			[]sorting.Sort{
				sorting.NewSort("created_at_desc", "Created date descending"),
			},
			"there's no order criteria with the id foo_desc",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var s sorting.Sorting
			err := sorting.NewDecoder(tc.urlParams, tc.criterias...).Decode(&s)
			assert.NotNil(t, err)
			assert.EqualError(t, err, tc.err)
		})
	}
}
