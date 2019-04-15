package middleware_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ifreddyrondon/bastion/middleware/listing/filtering"

	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/middleware/listing"
	"github.com/ifreddyrondon/bastion/middleware/listing/paging"
	"github.com/ifreddyrondon/bastion/middleware/listing/sorting"

	"github.com/ifreddyrondon/capture/pkg/http/rest/middleware"
)

func TestFilterUserRepos(t *testing.T) {
	t.Parallel()

	l := &listing.Listing{
		Paging: paging.Paging{
			Limit:           paging.DefaultLimit,
			Offset:          paging.DefaultOffset,
			MaxAllowedLimit: 50,
		},
		Sorting: &sorting.Sorting{
			Sort:      &updatedDESC,
			Available: []sorting.Sort{updatedDESC, updatedASC, createdDESC, createdASC},
		},
		Filtering: &filtering.Filtering{
			Available: []filtering.Filter{
				{
					ID:          "visibility",
					Description: "filters the repos by their visibility",
					Type:        "text",
					Values:      []filtering.Value{publicVisibility, privateVisibility},
				},
			},
		},
	}

	app, resultContainer := setupFilterMiddleware(middleware.FilterUserRepos())
	e := bastion.Tester(t, app)
	e.GET("/").
		Expect().
		Status(http.StatusOK)

	assert.Equal(t, l, resultContainer)
}

func TestFilterPublicRepos(t *testing.T) {
	t.Parallel()

	l := &listing.Listing{
		Paging: paging.Paging{
			Limit:           paging.DefaultLimit,
			Offset:          paging.DefaultOffset,
			MaxAllowedLimit: 50,
		},
		Sorting: &sorting.Sorting{
			Sort:      &updatedDESC,
			Available: []sorting.Sort{updatedDESC, updatedASC, createdDESC, createdASC},
		},
	}

	app, resultContainer := setupFilterMiddleware(middleware.FilterPublicRepos())
	e := bastion.Tester(t, app)
	e.GET("/").
		Expect().
		Status(http.StatusOK)

	assert.Equal(t, l, resultContainer)
}
