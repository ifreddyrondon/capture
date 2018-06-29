package listing

import (
	"net/http"

	"github.com/ifreddyrondon/bastion/render"
	"github.com/ifreddyrondon/bastion/render/json"
)

// Listing middleware
type Listing struct {
	defautls   *Params
	ctxManager *ContextManager
	render     render.Render
}

// NewListing retuns a new instance of Listing middleware.
// It receives a list of Option to modify the default values.
func NewListing(options ...Option) *Listing {
	l := &Listing{
		defautls:   NewParams(),
		ctxManager: NewContextManager(),
		render:     json.NewRender,
	}

	for _, o := range options {
		o(l)
	}

	return l
}

// List collect all the listing params and leaves them within the context of the request.
func (l *Listing) List(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var params Params
		if err := params.Decode(r.URL.Query(), *l.defautls); err != nil {
			_ = l.render(w).BadRequest(err)
			return
		}

		ctx := l.ctxManager.withParams(r.Context(), &params)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
