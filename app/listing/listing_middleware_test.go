package listing_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/ifreddyrondon/bastion/render/json"
	"github.com/ifreddyrondon/capture/app/listing"
	"github.com/ifreddyrondon/capture/app/listing/filtering"
	"github.com/ifreddyrondon/capture/app/listing/paging"
	"github.com/ifreddyrondon/capture/app/listing/sorting"
	"github.com/stretchr/testify/assert"
	httpexpect "gopkg.in/gavv/httpexpect.v1"
)

func setup(options ...listing.Option) (*httptest.Server, *listing.Listing, func()) {
	r := chi.NewRouter()
	var resultContainer listing.Listing

	midl := listing.NewParams(options...)
	r.Use(midl.Get)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		p, err := listing.NewContextManager().GetListing(r.Context())
		if err != nil {
			_ = json.NewRender(w).InternalServerError(err)
			return
		}

		resultContainer = *p
		w.Write([]byte("hi"))
	})

	server := httptest.NewServer(r)
	teardown := func() {
		server.Close()
	}
	return server, &resultContainer, teardown
}

func TestParamsMiddlewareFailure(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name      string
		urlParams string
		opts      []listing.Option
		response  map[string]interface{}
	}{
		{
			"given a bad offset param should return a 400",
			"offset=abc",
			[]listing.Option{listing.Limit(50)},
			map[string]interface{}{
				"status":  400.0,
				"error":   "Bad Request",
				"message": "invalid offset value, must be a number",
			},
		},
		{
			"given a sort query when none match sorting criteria should return a 400",
			"sort=foo_desc",
			[]listing.Option{listing.Sort(sorting.NewSort("created_at_desc", "Created date descending"))},
			map[string]interface{}{
				"status":  400.0,
				"error":   "Bad Request",
				"message": "there's no order criteria with the id foo_desc",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			server, _, teardown := setup(tc.opts...)
			defer teardown()
			e := httpexpect.New(t, server.URL)
			e.GET("/").WithQueryString(tc.urlParams).
				Expect().
				Status(http.StatusBadRequest).
				JSON().
				Object().Equal(tc.response)
		})
	}
}

func TestParamsMiddlewareOkWithOptions(t *testing.T) {
	t.Parallel()

	createdDescSort := sorting.NewSort("created_at_desc", "Created date descending")
	createdAscSort := sorting.NewSort("created_at_asc", "Created date ascendant")
	vNew := filtering.NewValue("new", "New")
	vUsed := filtering.NewValue("used", "Used")
	text := filtering.NewText("condition", "test", vNew, vUsed)
	vTrue := filtering.NewValue("true", "shared")
	vFalse := filtering.NewValue("false", "private")
	boolean := filtering.NewBoolean("shared", "test", "shared", "private")

	tt := []struct {
		name      string
		urlParams string
		opts      []listing.Option
		result    listing.Listing
	}{
		{
			"given non query params and not options should get default paging",
			"",
			[]listing.Option{},
			func() listing.Listing {
				return listing.Listing{
					Paging: paging.Paging{
						Limit:           paging.DefaultLimit,
						Offset:          paging.DefaultOffset,
						MaxAllowedLimit: paging.DefaultMaxAllowedLimit,
					},
				}
			}(),
		},
		{
			"given offset=11 params and not options should get paging with offset=11 and defaults",
			"offset=11",
			[]listing.Option{},
			func() listing.Listing {
				return listing.Listing{
					Paging: paging.Paging{
						Limit:           paging.DefaultLimit,
						Offset:          11,
						MaxAllowedLimit: paging.DefaultMaxAllowedLimit,
					},
				}
			}(),
		},
		{
			"given non query params and changing the default limit option should get default paging with limit 50",
			"",
			[]listing.Option{listing.Limit(50)},
			func() listing.Listing {
				return listing.Listing{
					Paging: paging.Paging{
						Limit:           50,
						Offset:          paging.DefaultOffset,
						MaxAllowedLimit: paging.DefaultMaxAllowedLimit,
					},
				}
			}(),
		},
		{
			"given limit=110 param and changing the default MaxAllowedLimit option to 120 should allow limit > 100 < 120",
			"limit=110",
			[]listing.Option{listing.MaxAllowedLimit(120)},
			func() listing.Listing {
				return listing.Listing{
					Paging: paging.Paging{
						Limit:           110,
						Offset:          paging.DefaultOffset,
						MaxAllowedLimit: 120,
					},
				}
			}(),
		},
		{
			"given non query params and one sort criteria should get sorting with default sort",
			"",
			[]listing.Option{listing.Sort(createdDescSort)},
			func() listing.Listing {
				return listing.Listing{
					Paging: paging.Paging{
						Limit:           paging.DefaultLimit,
						Offset:          paging.DefaultOffset,
						MaxAllowedLimit: paging.DefaultMaxAllowedLimit,
					},
					Sorting: &sorting.Sorting{
						Sort:      &createdDescSort,
						Available: []sorting.Sort{createdDescSort},
					},
				}
			}(),
		},
		{
			"given sort query params and sort criteria should get sorting with selected sort",
			"sort=created_at_desc",
			[]listing.Option{listing.Sort(createdDescSort, createdAscSort)},
			func() listing.Listing {
				return listing.Listing{
					Paging: paging.Paging{
						Limit:           paging.DefaultLimit,
						Offset:          paging.DefaultOffset,
						MaxAllowedLimit: paging.DefaultMaxAllowedLimit,
					},
					Sorting: &sorting.Sorting{
						Sort:      &createdDescSort,
						Available: []sorting.Sort{createdDescSort, createdAscSort},
					},
				}
			}(),
		},
		{
			"given non query params and one filter criteria should get filtering with only available",
			"",
			[]listing.Option{listing.Filter(text)},
			func() listing.Listing {
				return listing.Listing{
					Paging: paging.Paging{
						Limit:           paging.DefaultLimit,
						Offset:          paging.DefaultOffset,
						MaxAllowedLimit: paging.DefaultMaxAllowedLimit,
					},
					Filtering: &filtering.Filtering{
						Filters: []filtering.Filter{},
						Available: []filtering.Filter{
							filtering.Filter{
								ID:     "condition",
								Name:   "test",
								Type:   "text",
								Values: []filtering.Value{vNew, vUsed},
							},
						},
					},
				}
			}(),
		},
		{
			"given non query params and some filters criteria should get filtering with all available",
			"",
			[]listing.Option{listing.Filter(text, boolean)},
			func() listing.Listing {
				return listing.Listing{
					Paging: paging.Paging{
						Limit:           paging.DefaultLimit,
						Offset:          paging.DefaultOffset,
						MaxAllowedLimit: paging.DefaultMaxAllowedLimit,
					},
					Filtering: &filtering.Filtering{
						Filters: []filtering.Filter{},
						Available: []filtering.Filter{
							filtering.Filter{
								ID:     "condition",
								Name:   "test",
								Type:   "text",
								Values: []filtering.Value{vNew, vUsed},
							},
							filtering.Filter{
								ID:     "shared",
								Name:   "test",
								Type:   "boolean",
								Values: []filtering.Value{vTrue, vFalse},
							},
						},
					},
				}
			}(),
		},
		{
			"given a filter query params and some filters criteria should get filtering with all available and filter",
			"condition=new",
			[]listing.Option{listing.Filter(text, boolean)},
			func() listing.Listing {
				return listing.Listing{
					Paging: paging.Paging{
						Limit:           paging.DefaultLimit,
						Offset:          paging.DefaultOffset,
						MaxAllowedLimit: paging.DefaultMaxAllowedLimit,
					},
					Filtering: &filtering.Filtering{
						Filters: []filtering.Filter{
							filtering.Filter{
								ID:     "condition",
								Name:   "test",
								Type:   "text",
								Values: []filtering.Value{vNew},
							},
						},
						Available: []filtering.Filter{
							filtering.Filter{
								ID:     "condition",
								Name:   "test",
								Type:   "text",
								Values: []filtering.Value{vNew, vUsed},
							},
							filtering.Filter{
								ID:     "shared",
								Name:   "test",
								Type:   "boolean",
								Values: []filtering.Value{vTrue, vFalse},
							},
						},
					},
				}
			}(),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			server, resultContainer, teardown := setup(tc.opts...)
			defer teardown()
			e := httpexpect.New(t, server.URL)
			e.GET("/").WithQueryString(tc.urlParams).
				Expect().
				Status(http.StatusOK)

			assert.Equal(t, tc.result.Paging, resultContainer.Paging)
			if resultContainer.Sorting != nil {
				assert.Equal(t, tc.result.Sorting.Sort, resultContainer.Sorting.Sort)
				assert.Equal(t, tc.result.Sorting.Available, resultContainer.Sorting.Available)
			}
			if resultContainer.Filtering != nil {
				for i, f := range resultContainer.Filtering.Filters {
					assert.Equal(t, tc.result.Filtering.Filters[i], f)
				}
				assert.Equal(t, tc.result.Filtering.Available, resultContainer.Filtering.Available)
			}
		})
	}
}
