package repository

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/middleware"
	"github.com/ifreddyrondon/bastion/middleware/listing/sorting"
	"github.com/ifreddyrondon/bastion/render"
	"github.com/ifreddyrondon/capture/features"
	"github.com/ifreddyrondon/capture/features/repository/encoder"
)

// Routes returns a configured http.Handler with repositories resources.
func Routes(store Store) http.Handler {
	c := &controller{store: store, render: render.NewJSON()}

	updatedDESC := sorting.NewSort("updated_at_desc", "updated_at DESC", "Updated date descending")
	updatedASC := sorting.NewSort("updated_at_asc", "updated_at ASC", "Updated date ascendant")
	createdDESC := sorting.NewSort("created_at_desc", "created_at DESC", "Created date descending")
	createdASC := sorting.NewSort("created_at_asc", "created_at ASC", "Created date ascendant")

	listing := middleware.Listing(
		middleware.MaxAllowedLimit(50),
		middleware.Sort(updatedDESC, updatedASC, createdDESC, createdASC),
	)

	r := bastion.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Use(listing)
		r.Get("/", c.list)
	})

	return r
}

type controller struct {
	store  Store
	render render.APIRenderer
}

func (c *controller) list(w http.ResponseWriter, r *http.Request) {
	listing, err := middleware.GetListing(r.Context())
	if err != nil {
		c.render.InternalServerError(w, err)
		return
	}

	listingRepo := NewListingRepo(*listing)
	listingRepo.Visibility = &features.Public
	repos, err := c.store.List(listingRepo)
	if err != nil {
		c.render.InternalServerError(w, err)
		return
	}

	res := encoder.ListRepositoryResponse{Listing: listing, Results: repos}
	c.render.Send(w, res)
}
