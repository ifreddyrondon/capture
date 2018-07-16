package filtering_test

import (
	"testing"
)

func TestDecodeOK(t *testing.T) {
	t.Parallel()

	// textDecoder := filtering.NewTextDecoder()

	// tt := []struct {
	// 	name      string
	// 	urlParams url.Values
	// 	decoders  []filtering.FilterDecoder
	// 	result    filtering.Filtering
	// }{
	// {
	// 	"given none query params and non criterias should decode empty Sorting",
	// 	map[string][]string{},
	// 	[]sorting.Sort{},
	// 	sorting.Sorting{},
	// },
	// {
	// 	"given non sort query params present and one sort criteria",
	// 	map[string][]string{},
	// 	[]sorting.Sort{createdDescSort},
	// 	sorting.Sorting{
	// 		Sort:      &createdDescSort,
	// 		Available: []sorting.Sort{createdDescSort},
	// 	},
	// },
	// {
	// 	"given non sort query params present and one some sort criteria",
	// 	map[string][]string{},
	// 	[]sorting.Sort{createdDescSort, createdAscSort},
	// 	sorting.Sorting{
	// 		Sort:      &createdDescSort,
	// 		Available: []sorting.Sort{createdDescSort, createdAscSort},
	// 	},
	// },
	// {
	// 	"given created_at_desc sort query params present and one some sort criteria",
	// 	map[string][]string{"sort": []string{"created_at_asc"}},
	// 	[]sorting.Sort{createdDescSort, createdAscSort},
	// 	sorting.Sorting{
	// 		Sort:      &createdAscSort,
	// 		Available: []sorting.Sort{createdDescSort, createdAscSort},
	// 	},
	// },
	// }

	// for _, tc := range tt {
	// 	t.Run(tc.name, func(t *testing.T) {
	// 		var f filtering.Filtering
	// 		err := filtering.NewDecoder(tc.urlParams, tc.decoders...).Decode(&f)
	// 		assert.Nil(t, err)
	// 		assert.Equal(t, tc.result.Filters, f.Filters)
	// 		assert.Equal(t, len(tc.result.Available), len(f.Available))
	// 	})
	// }
}
