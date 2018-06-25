package listing_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ifreddyrondon/capture/app/listing"

	"github.com/go-chi/chi"
	httpexpect "gopkg.in/gavv/httpexpect.v1"
)

func TestA(t *testing.T) {
	r := chi.NewRouter()

	listingMiddl := listing.NewListing()
	r.Use(listingMiddl.List)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hi"))
	})

	server := httptest.NewServer(r)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	// is it working?
	e.GET("/").WithQuery("offset", "11").
		Expect().
		Status(http.StatusOK)
}
