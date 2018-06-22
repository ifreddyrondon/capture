package search_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/ifreddyrondon/capture/app/search"
	httpexpect "gopkg.in/gavv/httpexpect.v1"
)

func TestA(t *testing.T) {
	r := chi.NewRouter()

	r.Use(search.SearchParams)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hi"))
	})

	server := httptest.NewServer(r)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	// is it working?
	e.GET("/").WithQuery("limit", "1").WithQuery("offset", 10).
		Expect().
		Status(http.StatusOK)
}
