package listing_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/ifreddyrondon/bastion/render/json"
	"github.com/ifreddyrondon/capture/app/listing"
	"github.com/ifreddyrondon/capture/app/listing/paging"
	"github.com/stretchr/testify/assert"
	httpexpect "gopkg.in/gavv/httpexpect.v1"
)

func setup(options ...listing.Option) (*httptest.Server, *listing.Listing, func()) {
	r := chi.NewRouter()
	var resultContainer listing.Listing

	midl := listing.NewURLParams(options...)
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

func TestListingMiddlewareFailure(t *testing.T) {
	t.Parallel()

	response := map[string]interface{}{
		"status":  400.0,
		"error":   "Bad Request",
		"message": "invalid offset value, must be a number",
	}

	server, _, teardown := setup()
	defer teardown()
	e := httpexpect.New(t, server.URL)
	e.GET("/").WithQuery("offset", "abc").
		Expect().
		Status(http.StatusBadRequest).
		JSON().
		Object().Equal(response)
}

func TestListingMiddlewareGettingDefaults(t *testing.T) {
	t.Parallel()

	server, resultContainer, teardown := setup()
	defer teardown()
	e := httpexpect.New(t, server.URL)
	e.GET("/").
		Expect().
		Status(http.StatusOK)

	assert.Equal(t, 10, resultContainer.Paging.Limit)
	assert.Equal(t, 100, resultContainer.Paging.MaxAllowedLimit)
	assert.Equal(t, int64(0), resultContainer.Paging.Offset)
}

func TestListingMiddlewareOkGettingOffset(t *testing.T) {
	t.Parallel()

	server, resultContainer, teardown := setup()
	defer teardown()
	e := httpexpect.New(t, server.URL)
	e.GET("/").WithQuery("offset", "11").
		Expect().
		Status(http.StatusOK)

	assert.Equal(t, 10, resultContainer.Paging.Limit)
	assert.Equal(t, int64(11), resultContainer.Paging.Offset)
}

func TestListingMiddlewareOkWithOptions(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name      string
		urlParams string
		opts      []listing.Option
		result    listing.Listing
	}{
		{
			"get new default limit changed by option",
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
			"get limit by query when change max allowed limit",
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
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			server, resultContainer, teardown := setup(tc.opts...)
			defer teardown()
			e := httpexpect.New(t, server.URL)
			e.GET("/").WithQueryString(tc.urlParams).
				Expect().
				Status(http.StatusOK)

			assert.Equal(t, tc.result.Paging.Limit, resultContainer.Paging.Limit)
			assert.Equal(t, tc.result.Paging.Offset, resultContainer.Paging.Offset)
			assert.Equal(t, tc.result.Paging.MaxAllowedLimit, resultContainer.Paging.MaxAllowedLimit)
		})
	}
}
