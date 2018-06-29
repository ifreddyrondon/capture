package listing_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ifreddyrondon/capture/app/listing"
	"github.com/stretchr/testify/assert"

	"github.com/go-chi/chi"
	"github.com/ifreddyrondon/bastion/render/json"
	httpexpect "gopkg.in/gavv/httpexpect.v1"
)

func setup(options ...listing.Option) (*httptest.Server, *listing.Params, func()) {
	r := chi.NewRouter()
	var resultContainer listing.Params

	listingMiddl := listing.NewListing(options...)
	r.Use(listingMiddl.List)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		p, err := listing.NewContextManager().GetParams(r.Context())
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
	assert.Equal(t, int64(0), resultContainer.Paging.Offset)
}

func TestListingMiddlewareGettingDefaultsDefinedByOptions(t *testing.T) {
	t.Parallel()

	server, resultContainer, teardown := setup(listing.Limit(100))
	defer teardown()
	e := httpexpect.New(t, server.URL)
	e.GET("/").
		Expect().
		Status(http.StatusOK)

	assert.Equal(t, 100, resultContainer.Paging.Limit)
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

func TestListingMiddlewareOkGettingLimit(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name       string
		limitQuery string
		opts       []listing.Option
		result     listing.Paging
	}{
		{
			"get limit by query",
			"11",
			[]listing.Option{},
			func() listing.Paging {
				p := listing.NewPaging()
				p.Limit = 11
				return p
			}(),
		},
		{
			"get limit by query when change max allowed limit ",
			"110",
			[]listing.Option{listing.MaxAllowedLimit(120)},
			func() listing.Paging {
				p := listing.NewPaging()
				p.Limit = 110
				return p
			}(),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			server, resultContainer, teardown := setup(tc.opts...)
			defer teardown()
			e := httpexpect.New(t, server.URL)
			e.GET("/").WithQuery("limit", tc.limitQuery).
				Expect().
				Status(http.StatusOK)

			assert.Equal(t, tc.result.Limit, resultContainer.Paging.Limit)
			assert.Equal(t, tc.result.Offset, resultContainer.Paging.Offset)
		})
	}
}
