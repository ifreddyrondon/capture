package domain

import (
	"github.com/ifreddyrondon/bastion/middleware/listing"
	"github.com/ifreddyrondon/capture/pkg"
)

// Listing allows to sort and filter search results between services and storage.
type Listing struct {
	SortKey    string
	Offset     int64
	Limit      int
	Owner      string
	Visibility *pkg.Visibility
}

// NewListing returns a new Listing instance with offset and limit from listing.Listing.
// It'll get SortKey and Visibility from Sorting and Filtering if there are available.
func NewListing(l listing.Listing) *Listing {
	domainListing := &Listing{
		Offset: l.Paging.Offset,
		Limit:  l.Paging.Limit,
	}

	if l.Sorting != nil {
		domainListing.SortKey = l.Sorting.Sort.Value
	}

	if l.Filtering == nil {
		return domainListing
	}

	for i := range l.Filtering.Filters {
		if l.Filtering.Filters[i].ID == "visibility" {
			visibility := pkg.Visibility(l.Filtering.Filters[i].Values[0].ID)
			domainListing.Visibility = &visibility
			break
		}
	}

	return domainListing
}
