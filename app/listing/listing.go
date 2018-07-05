package listing

import (
	"github.com/ifreddyrondon/capture/app/listing/paging"
	"github.com/ifreddyrondon/capture/app/listing/sorting"
)

// Listing containst the info to perform filter sort and paging over a collection.
type Listing struct {
	Paging paging.Paging
	sorting.Sorting
	// AvailableFilter []Filter
	// Filter          Filter
}

// // FilterValue defines a value that a Filter can have.
// type FilterValue struct {
// 	ID, Name string
// }

// // Filter struct allows to filter a collection by an identifier.
// type Filter struct {
// 	ID, Name string
// 	values   []FilterValue
// }
