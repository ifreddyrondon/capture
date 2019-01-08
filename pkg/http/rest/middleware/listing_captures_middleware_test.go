package middleware_test

import (
	"net/http"
	"testing"

	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/middleware/listing"
	"github.com/ifreddyrondon/bastion/middleware/listing/paging"
	"github.com/ifreddyrondon/bastion/middleware/listing/sorting"
	"github.com/ifreddyrondon/capture/pkg/http/rest/middleware"
	"github.com/stretchr/testify/assert"
)

func TestFilterCaptures(t *testing.T) {
	t.Parallel()

	l := &listing.Listing{
		Paging: paging.Paging{
			Limit:           paging.DefaultLimit,
			Offset:          paging.DefaultOffset,
			MaxAllowedLimit: 100,
		},
		Sorting: &sorting.Sorting{
			Sort:      &updatedDESC,
			Available: []sorting.Sort{updatedDESC, updatedASC, createdDESC, createdASC},
		},
	}

	app, resultContainer := setupFilterMiddleware(middleware.FilterCaptures())
	e := bastion.Tester(t, app)
	e.GET("/").
		Expect().
		Status(http.StatusOK)

	assert.Equal(t, l, resultContainer)
}
