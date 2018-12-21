package rest_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/ifreddyrondon/bastion"
	bastionMiddleware "github.com/ifreddyrondon/bastion/middleware"
	listingBastionMiddleware "github.com/ifreddyrondon/bastion/middleware/listing"
	"github.com/ifreddyrondon/bastion/middleware/listing/filtering"
	"github.com/ifreddyrondon/bastion/middleware/listing/paging"
	"github.com/ifreddyrondon/bastion/middleware/listing/sorting"
	"github.com/ifreddyrondon/capture/pkg"
	"github.com/ifreddyrondon/capture/pkg/domain"
	"github.com/ifreddyrondon/capture/pkg/http/rest"
	"github.com/ifreddyrondon/capture/pkg/listing"
)

func listingMiddlewareOK(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		createdDESC := sorting.NewSort("created_at_desc", "created_at DESC", "Created date descending")
		createdASC := sorting.NewSort("created_at_asc", "created_at ASC", "Created date ascendant")
		publicVisibility := filtering.NewValue("public", "public repos")
		privateVisibility := filtering.NewValue("private", "private repos")

		l := &listingBastionMiddleware.Listing{
			Paging: paging.Paging{
				Limit:  50,
				Offset: 0,
			},
			Sorting: &sorting.Sorting{
				Sort:      &createdDESC,
				Available: []sorting.Sort{createdDESC, createdASC},
			},
			Filtering: &filtering.Filtering{
				Filters: []filtering.Filter{
					{
						ID:          "visibility",
						Description: "test",
						Type:        "text",
						Values:      []filtering.Value{publicVisibility},
					},
				},
				Available: []filtering.Filter{
					{
						ID:          "visibility",
						Description: "test",
						Type:        "text",
						Values:      []filtering.Value{publicVisibility, privateVisibility},
					},
				},
			},
		}

		ctx := context.WithValue(r.Context(), bastionMiddleware.ListingCtxKey, l)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func listingMiddlewareBAD(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), bastionMiddleware.ListingCtxKey, "bad listing")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type mockListingService struct {
	repos []pkg.Repository
	err   error
}

func (m *mockListingService) GetUserRepos(u *domain.User, l *listingBastionMiddleware.Listing) (*listing.ListRepositoryResponse, error) {
	return &listing.ListRepositoryResponse{Listing: l, Results: m.repos}, m.err
}

func (m *mockListingService) GetPublicRepos(l *listingBastionMiddleware.Listing) (*listing.ListRepositoryResponse, error) {
	return &listing.ListRepositoryResponse{Listing: l, Results: m.repos}, m.err
}

func setupListingPublicReposHandler(s listing.Service, listMiddle func(http.Handler) http.Handler) *bastion.Bastion {
	app := bastion.New()
	app.APIRouter.Use(listMiddle)
	app.APIRouter.Get("/", rest.ListingPublicRepos(s))
	return app
}

func setupListingUserReposHandler(s listing.Service, listMiddle, auth func(http.Handler) http.Handler) *bastion.Bastion {
	app := bastion.New()
	app.APIRouter.Use(listMiddle)
	app.APIRouter.Use(auth)
	app.APIRouter.Get("/", rest.ListingUserRepos(s))
	return app
}

func TestListingPublicReposSuccess(t *testing.T) {
	t.Parallel()

	repos := []pkg.Repository{
		{Name: "test public", Visibility: "public"},
		{Name: "test private", Visibility: "private"},
	}
	s := &mockListingService{repos: repos}
	app := setupListingPublicReposHandler(s, listingMiddlewareOK)

	e := bastion.Tester(t, app)
	res := e.GET("/").Expect().
		JSON().Object()
	res.Value("results").Array().Length().Equal(2)
	res.Value("listing").Object().
		ContainsKey("paging").
		ContainsKey("sorting").
		ContainsKey("filtering")
}

func TestListingPublicReposInternalServerBadListing(t *testing.T) {
	t.Parallel()

	s := &mockListingService{}
	app := setupListingPublicReposHandler(s, listingMiddlewareBAD)

	response := map[string]interface{}{
		"status":  500.0,
		"error":   "Internal Server Error",
		"message": "looks like something went wrong",
	}

	e := bastion.Tester(t, app)
	e.GET("/").Expect().
		Status(http.StatusInternalServerError).
		JSON().Object().Equal(response)
}

func TestListingPublicReposInternalServerGettingRepos(t *testing.T) {
	t.Parallel()

	s := &mockListingService{err: errors.New("test")}
	app := setupListingPublicReposHandler(s, listingMiddlewareOK)

	response := map[string]interface{}{
		"status":  500.0,
		"error":   "Internal Server Error",
		"message": "looks like something went wrong",
	}

	e := bastion.Tester(t, app)
	e.GET("/").Expect().
		Status(http.StatusInternalServerError).
		JSON().Object().Equal(response)
}

func TestListingUserReposSuccess(t *testing.T) {
	t.Parallel()

	repos := []pkg.Repository{
		{Name: "test public", Visibility: "public"},
		{Name: "test private", Visibility: "private"},
	}
	s := &mockListingService{repos: repos}
	app := setupListingUserReposHandler(s, listingMiddlewareOK, authenticatedMiddleware)

	e := bastion.Tester(t, app)
	res := e.GET("/").Expect().
		JSON().Object()
	res.Value("results").Array().Length().Equal(2)
	res.Value("listing").Object().
		ContainsKey("paging").
		ContainsKey("sorting").
		ContainsKey("filtering")

	results := res.Value("results").Array()
	results.First().Object().
		ValueEqual("name", "test public").
		ValueEqual("visibility", "public")
}

func TestListingUserReposInternalServerBadListing(t *testing.T) {
	t.Parallel()

	s := &mockListingService{err: errors.New("test")}
	app := setupListingUserReposHandler(s, listingMiddlewareBAD, authenticatedMiddleware)

	response := map[string]interface{}{
		"status":  500.0,
		"error":   "Internal Server Error",
		"message": "looks like something went wrong",
	}

	e := bastion.Tester(t, app)
	e.GET("/").Expect().
		Status(http.StatusInternalServerError).
		JSON().Object().Equal(response)
}

func TestListingUserReposFailInternalErrorGettingUser(t *testing.T) {
	t.Parallel()
	s := &mockListingService{}
	app := setupListingUserReposHandler(s, listingMiddlewareOK, notUserMiddleware)

	response := map[string]interface{}{
		"status":  500.0,
		"error":   "Internal Server Error",
		"message": "looks like something went wrong",
	}

	e := bastion.Tester(t, app)
	e.GET("/").Expect().
		Status(http.StatusInternalServerError).
		JSON().Object().Equal(response)
}

func TestListingUserReposInternalServerGettingRepos(t *testing.T) {
	t.Parallel()

	s := &mockListingService{err: errors.New("test")}
	app := setupListingUserReposHandler(s, listingMiddlewareOK, authenticatedMiddleware)

	response := map[string]interface{}{
		"status":  500.0,
		"error":   "Internal Server Error",
		"message": "looks like something went wrong",
	}

	e := bastion.Tester(t, app)
	e.GET("/").Expect().
		Status(http.StatusInternalServerError).
		JSON().Object().Equal(response)
}
