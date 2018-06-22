package search

import (
	"fmt"
	"net/http"

	"github.com/gorilla/schema"
)

type name struct {
}

func SearchParams(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// params := r.URL.Query()
		// fmt.Println(params)
		if err := r.ParseForm(); err != nil {
			// Handle error
		}

		paging := new(Paging)
		if err := schema.NewDecoder().Decode(paging, r.Form); err != nil {
			// Handle error
		}

		fmt.Printf("%+v", paging)

		next.ServeHTTP(w, r)
	})
}
