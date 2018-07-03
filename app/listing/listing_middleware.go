package listing

import (
	"net/http"

	"github.com/ifreddyrondon/bastion/render"
	"github.com/ifreddyrondon/bastion/render/json"
)

// URLParams is a middleware that parses the url query params from a request and stores it
// on the context as a Listing under the key `listing_value`. It can be accessed through
// listing.GetListing.
//
// Sample usage.. for the url: `/repositories/1?limit=10&offset=25`
//
//  func routes() http.Handler {
//    r := chi.NewRouter()
//    r.Use(listing.URLListing)
//
//    r.Get("/repositories/{id}", ListRepositories)
//
//    return r
//  }
//
//  func ListRepositories(w http.ResponseWriter, r *http.Request) {
// 	  list, _ := listing.GetListing(r.Context())
//
// 	  // do something with listing
// }
type URLParams struct {
	optionsDecoder []func(*Decoder)
	ctxManager     *ContextManager
	render         render.Render
}

// Option allows to modify the defaults midleware values.
type Option func(*URLParams)

// Limit set the paging limit default.
func Limit(limit int) Option {
	return func(dec *URLParams) {
		o := DecodeLimit(limit)
		dec.optionsDecoder = append(dec.optionsDecoder, o)
	}
}

// MaxAllowedLimit set the max allowed limit default.
func MaxAllowedLimit(maxAllowed int) Option {
	return func(dec *URLParams) {
		o := DecodeMaxAllowedLimit(maxAllowed)
		dec.optionsDecoder = append(dec.optionsDecoder, o)
	}
}

// NewURLParams retuns a new instance of Listing middleware.
// It receives a list of Option to modify the default values.
func NewURLParams(options ...Option) *URLParams {
	l := &URLParams{
		ctxManager: NewContextManager(),
		render:     json.NewRender,
	}

	for _, o := range options {
		o(l)
	}

	return l
}

// Get collect all the listing params and stores it on the context
// as a Listing under the key `listing_value`.
func (m *URLParams) Get(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var l Listing

		if err := NewDecoder(r.URL.Query(), m.optionsDecoder...).Decode(&l); err != nil {
			_ = m.render(w).BadRequest(err)
			return
		}

		ctx := m.ctxManager.withParams(r.Context(), &l)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
