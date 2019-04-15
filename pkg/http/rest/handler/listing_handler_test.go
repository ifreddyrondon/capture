package handler_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/ifreddyrondon/capture/pkg/http/rest/handler"

	"github.com/ifreddyrondon/bastion"
	bastionMiddleware "github.com/ifreddyrondon/bastion/middleware"
	listingBastionMiddleware "github.com/ifreddyrondon/bastion/middleware/listing"
	"github.com/ifreddyrondon/bastion/middleware/listing/filtering"
	"github.com/ifreddyrondon/bastion/middleware/listing/paging"
	"github.com/ifreddyrondon/bastion/middleware/listing/sorting"
	"gopkg.in/src-d/go-kallax.v1"

	"github.com/ifreddyrondon/capture/pkg/domain"
	"github.com/ifreddyrondon/capture/pkg/listing"
)

func getBaseListing() listingBastionMiddleware.Listing {
	createdDESC := sorting.NewSort("created_at_desc", "created_at DESC", "Created date descending")
	createdASC := sorting.NewSort("created_at_asc", "created_at ASC", "Created date ascendant")
	return listingBastionMiddleware.Listing{
		Paging: paging.Paging{Limit: 50, Offset: 0},
		Sorting: &sorting.Sorting{
			Sort:      &createdDESC,
			Available: []sorting.Sort{createdDESC, createdASC},
		},
	}
}

func listingRepoMiddlewareOK(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		publicVisibility := filtering.NewValue("public", "public repos")
		privateVisibility := filtering.NewValue("private", "private repos")

		l := getBaseListing()
		l.Filtering = &filtering.Filtering{
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
		}

		ctx := context.WithValue(r.Context(), bastionMiddleware.ListingCtxKey, &l)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func listingMiddlewareBAD(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), bastionMiddleware.ListingCtxKey, "bad listing")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type mockListingRepoService struct {
	repos []domain.Repository
	err   error
}

func (m *mockListingRepoService) GetUserRepos(u *domain.User, l *listingBastionMiddleware.Listing) (*listing.ListRepositoryResponse, error) {
	return &listing.ListRepositoryResponse{Listing: l, Results: m.repos}, m.err
}

func (m *mockListingRepoService) GetPublicRepos(l *listingBastionMiddleware.Listing) (*listing.ListRepositoryResponse, error) {
	return &listing.ListRepositoryResponse{Listing: l, Results: m.repos}, m.err
}

func setupListingPublicReposHandler(s listing.RepoService, listMiddle func(http.Handler) http.Handler) *bastion.Bastion {
	app := bastion.New()
	app.Use(listMiddle)
	app.Get("/", handler.ListingPublicRepos(s))
	return app
}

func setupListingUserReposHandler(s listing.RepoService, listMiddle, auth func(http.Handler) http.Handler) *bastion.Bastion {
	app := bastion.New()
	app.Use(listMiddle)
	app.Use(auth)
	app.Get("/", handler.ListingUserRepos(s))
	return app
}

func TestListingPublicReposSuccess(t *testing.T) {
	t.Parallel()

	repos := []domain.Repository{
		{Name: "test public", Visibility: "public"},
		{Name: "test private", Visibility: "private"},
	}
	s := &mockListingRepoService{repos: repos}
	app := setupListingPublicReposHandler(s, listingRepoMiddlewareOK)

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

	s := &mockListingRepoService{}
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

	s := &mockListingRepoService{err: errors.New("test")}
	app := setupListingPublicReposHandler(s, listingRepoMiddlewareOK)

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

	repos := []domain.Repository{
		{Name: "test public", Visibility: "public"},
		{Name: "test private", Visibility: "private"},
	}
	s := &mockListingRepoService{repos: repos}
	app := setupListingUserReposHandler(s, listingRepoMiddlewareOK, withUserMiddle(defaultUser))

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

	s := &mockListingRepoService{err: errors.New("test")}
	app := setupListingUserReposHandler(s, listingMiddlewareBAD, withUserMiddle(defaultUser))

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
	s := &mockListingRepoService{}
	app := setupListingUserReposHandler(s, listingRepoMiddlewareOK, withUserMiddle(nil))

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

	s := &mockListingRepoService{err: errors.New("test")}
	app := setupListingUserReposHandler(s, listingRepoMiddlewareOK, withUserMiddle(defaultUser))

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

func listingCaptureMiddlewareOK(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := getBaseListing()
		ctx := context.WithValue(r.Context(), bastionMiddleware.ListingCtxKey, &l)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type mockListingCaptureService struct {
	captures []domain.Capture
	err      error
}

func (m *mockListingCaptureService) ListRepoCaptures(r *domain.Repository, l *listingBastionMiddleware.Listing) (*listing.ListCaptureResponse, error) {
	return &listing.ListCaptureResponse{Listing: l, Results: m.captures}, m.err
}

func setupListingRepoCapturesHandler(s listing.CaptureService, listMiddle, auth, repoMiddle func(http.Handler) http.Handler) *bastion.Bastion {
	app := bastion.New()
	app.Use(listMiddle)
	app.Use(auth)
	app.Use(repoMiddle)
	app.Get("/", handler.ListingRepoCaptures(s))
	return app
}

func TestListingRepoCapturesSuccess(t *testing.T) {
	t.Parallel()

	captures := []domain.Capture{{ID: kallax.NewULID()}, {ID: kallax.NewULID()}}
	s := &mockListingCaptureService{captures: captures}
	app := setupListingRepoCapturesHandler(s, listingCaptureMiddlewareOK, withUserMiddle(defaultUser), withRepoMiddle(defaultRepo))

	e := bastion.Tester(t, app)
	res := e.GET("/").Expect().
		JSON().Object()
	res.Value("results").Array().Length().Equal(2)
	res.Value("listing").Object().
		ContainsKey("paging").
		ContainsKey("sorting").
		NotContainsKey("filtering")
}

func TestListingRepoCapturesInternalServerBadListing(t *testing.T) {
	t.Parallel()

	s := &mockListingCaptureService{err: errors.New("test")}
	app := setupListingRepoCapturesHandler(s, listingMiddlewareBAD, withUserMiddle(defaultUser), withRepoMiddle(defaultRepo))

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

func TestListingRepoCapturesFailInternalErrorGettingRepo(t *testing.T) {
	t.Parallel()
	s := &mockListingCaptureService{}
	app := setupListingRepoCapturesHandler(s, listingRepoMiddlewareOK, withUserMiddle(defaultUser), withRepoMiddle(nil))

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

func TestListingUserReposInternalServerGettingCaptures(t *testing.T) {
	t.Parallel()

	s := &mockListingCaptureService{err: errors.New("test")}
	app := setupListingRepoCapturesHandler(s, listingCaptureMiddlewareOK, withUserMiddle(defaultUser), withRepoMiddle(defaultRepo))

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
