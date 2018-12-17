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

// NewListing returns a new Listing instance with SortKey, offset and limit from listing.Listing.
func NewListing(l listing.Listing) *Listing {
	domainListing := &Listing{
		SortKey: l.Sorting.Sort.Value,
		Offset:  l.Paging.Offset,
		Limit:   l.Paging.Limit,
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
