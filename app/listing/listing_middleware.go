package listing

import (
	"fmt"
	"net/http"

	"github.com/ifreddyrondon/capture/app/listing/paging"
)

// Params containst the info to perform filter sort and paging over a collection.
type Params struct {
	paging.Paging
	AvailableSort []Sort
	Sort
	AvailableFilter []Filter
	Filter          Filter
}

func defaultParams() *Params {
	return &Params{
		Paging: paging.NewDefaults(),
	}
}

// Listing middleware
type Listing struct {
	defautls *Params
}

// Options are function to modify the defaults params values
type Options func(*Listing)

// NewListing retuns a new instance of Listing middleware.
// It receives a list of options to modify the default values.
func NewListing(...Options) *Listing {
	l := &Listing{defautls: defaultParams()}

	return l
}

// List collect all the listing params and leaves them within the context of the request.
func (l *Listing) List(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		pag := &paging.Paging{}
		pag.Decode(params, l.defautls.Paging)
		fmt.Printf("%+v\n", pag)

		next.ServeHTTP(w, r)
	})
}
