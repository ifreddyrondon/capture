package listing

import (
	"fmt"
	"net/http"
)

// Listing middleware
type Listing struct {
	defautls *Params
}

// Options are function to modify the defaults params values
type Options func(*Listing)

// NewListing retuns a new instance of Listing middleware.
// It receives a list of options to modify the default values.
func NewListing(...Options) *Listing {
	l := &Listing{defautls: NewParamsDefault()}

	return l
}

// List collect all the listing params and leaves them within the context of the request.
func (l *Listing) List(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var params Params
		params.Decode(r.URL.Query(), *l.defautls)
		fmt.Printf("%+v\n", params)

		next.ServeHTTP(w, r)
	})
}
