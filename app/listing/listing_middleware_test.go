package listing_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ifreddyrondon/capture/app/listing"

	"github.com/stretchr/testify/assert"
	httpexpect "gopkg.in/gavv/httpexpect.v1"

	"github.com/ifreddyrondon/bastion/render/json"

	"github.com/go-chi/chi"
)

func setup() (*httptest.Server, *listing.Params, func()) {
	r := chi.NewRouter()
	var resultContainer listing.Params

	listingMiddl := listing.NewListing()
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

func TestListingMiddlewareOkWithOffset(t *testing.T) {
	t.Parallel()

	server, resultContainer, teardown := setup()
	defer teardown()
	e := httpexpect.New(t, server.URL)
	e.GET("/").WithQuery("offset", "11").
		Expect().
		Status(http.StatusOK)

	assert.Equal(t, int64(10), resultContainer.Paging.Limit)
	assert.Equal(t, int64(11), resultContainer.Paging.Offset)
}

func TestListingMiddlewareOkWithLimit(t *testing.T) {
	t.Parallel()

	server, resultContainer, teardown := setup()
	defer teardown()
	e := httpexpect.New(t, server.URL)
	e.GET("/").WithQuery("limit", "11").
		Expect().
		Status(http.StatusOK)

	assert.Equal(t, int64(11), resultContainer.Paging.Limit)
	assert.Equal(t, int64(0), resultContainer.Paging.Offset)
}

func TestListingMiddlewareFailure(t *testing.T) {
	t.Parallel()

	response := map[string]interface{}{
		"status":  400.0,
		"error":   "Bad Request",
		"message": "invalid offset value",
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
