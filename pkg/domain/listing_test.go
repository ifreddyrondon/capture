package domain_test

import (
	"testing"

	"github.com/ifreddyrondon/bastion/middleware/listing"
	"github.com/ifreddyrondon/bastion/middleware/listing/filtering"
	"github.com/ifreddyrondon/bastion/middleware/listing/paging"
	"github.com/ifreddyrondon/bastion/middleware/listing/sorting"
	"github.com/stretchr/testify/assert"

	"github.com/ifreddyrondon/capture/pkg/domain"
)

func TestNewListingWithoutVisibilityFilterAndSorting(t *testing.T) {
	t.Parallel()
	l := listing.Listing{
		Paging: paging.Paging{
			Limit:  50,
			Offset: 0,
		},
	}

	result := domain.NewListing(l)
	assert.Equal(t, 50, result.Limit)
	assert.Equal(t, int64(0), result.Offset)
	assert.Equal(t, "", result.SortKey)
	assert.Nil(t, result.Owner)
	assert.Nil(t, result.Visibility)
}

func TestNewListingWithSortingAndWithoutVisibilityFilter(t *testing.T) {
	t.Parallel()
	createdDESC := sorting.NewSort("created_at_desc", "created_at DESC", "Created date descending")
	createdASC := sorting.NewSort("created_at_asc", "created_at ASC", "Created date ascendant")

	l := listing.Listing{
		Paging: paging.Paging{
			Limit:  50,
			Offset: 0,
		},
		Sorting: &sorting.Sorting{
			Sort:      &createdDESC,
			Available: []sorting.Sort{createdDESC, createdASC},
		},
	}

	result := domain.NewListing(l)
	assert.Equal(t, 50, result.Limit)
	assert.Equal(t, int64(0), result.Offset)
	assert.Equal(t, "created_at DESC", result.SortKey)
	assert.Nil(t, result.Owner)
	assert.Nil(t, result.Visibility)
}

func TestNewListingWithVisibilityFilterNotAppliedAndWithoutSorting(t *testing.T) {
	t.Parallel()
	publicVisibility := filtering.NewValue("public", "public repos")
	privateVisibility := filtering.NewValue("private", "private repos")

	l := listing.Listing{
		Paging: paging.Paging{
			Limit:  50,
			Offset: 0,
		},
		Filtering: &filtering.Filtering{
			Available: []filtering.Filter{
				{
					ID:          "visibility",
					Description: "test",
					Type:        "text",
					Values:      []filtering.Value{publicVisibility, privateVisibility},
				},
			},
		},
	}

	result := domain.NewListing(l)
	assert.Equal(t, 50, result.Limit)
	assert.Equal(t, int64(0), result.Offset)
	assert.Equal(t, "", result.SortKey)
	assert.Nil(t, result.Owner)
	assert.Nil(t, result.Visibility)
}

func TestNewListingWithVisibilityFilterAppliedAndWithoutSorting(t *testing.T) {
	t.Parallel()
	publicVisibility := filtering.NewValue("public", "public repos")
	privateVisibility := filtering.NewValue("private", "private repos")

	l := listing.Listing{
		Paging: paging.Paging{
			Limit:  50,
			Offset: 0,
		},
		Filtering: &filtering.Filtering{
			Filters: []filtering.Filter{
				{
					ID:          "visibility",
					Description: "test",
					Type:        "text",
					Values:      []filtering.Value{publicVisibility},
				},
			},
			Available: []filtering.Filter{
				{
					ID:          "visibility",
					Description: "test",
					Type:        "text",
					Values:      []filtering.Value{publicVisibility, privateVisibility},
				},
			},
		},
	}

	result := domain.NewListing(l)
	assert.Equal(t, 50, result.Limit)
	assert.Equal(t, int64(0), result.Offset)
	assert.Equal(t, "", result.SortKey)
	assert.Nil(t, result.Owner)
	assert.Equal(t, &domain.Public, result.Visibility)
}

func TestNewListingWithVisibilityFilterAppliedAndSorting(t *testing.T) {
	t.Parallel()
	createdDESC := sorting.NewSort("created_at_desc", "created_at DESC", "Created date descending")
	createdASC := sorting.NewSort("created_at_asc", "created_at ASC", "Created date ascendant")
	publicVisibility := filtering.NewValue("public", "public repos")
	privateVisibility := filtering.NewValue("private", "private repos")

	l := listing.Listing{
		Paging: paging.Paging{
			Limit:  50,
			Offset: 0,
		},
		Sorting: &sorting.Sorting{
			Sort:      &createdDESC,
			Available: []sorting.Sort{createdDESC, createdASC},
		},
		Filtering: &filtering.Filtering{
			Filters: []filtering.Filter{
				{
					ID:          "visibility",
					Description: "test",
					Type:        "text",
					Values:      []filtering.Value{publicVisibility},
				},
			},
			Available: []filtering.Filter{
				{
					ID:          "visibility",
					Description: "test",
					Type:        "text",
					Values:      []filtering.Value{publicVisibility, privateVisibility},
				},
			},
		},
	}

	result := domain.NewListing(l)
	assert.Equal(t, 50, result.Limit)
	assert.Equal(t, int64(0), result.Offset)
	assert.Equal(t, "created_at DESC", result.SortKey)
	assert.Nil(t, result.Owner)
	assert.Equal(t, &domain.Public, result.Visibility)
}
