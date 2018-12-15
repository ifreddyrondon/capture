package repository_test

import (
	"net/http"
	"testing"

	"github.com/ifreddyrondon/bastion/middleware"

	"github.com/ifreddyrondon/bastion/middleware/listing/filtering"
	"github.com/ifreddyrondon/bastion/middleware/listing/sorting"

	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/capture/pkg"
	auth "github.com/ifreddyrondon/capture/pkg/http/rest/middleware"
	"github.com/ifreddyrondon/capture/pkg/repository"
)

var tempUser = pkg.User{Email: "test@example.com", ID: "0162eb39-a65e-04a1-7ad9-d663bb49a396"}

func authRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := auth.WithUser(r.Context(), &tempUser)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type mockStore struct {
	repo  *pkg.Repository
	repos []pkg.Repository
	err   error
}

func (m *mockStore) Drop() {}
func (m *mockStore) List(l repository.ListingRepo) ([]pkg.Repository, error) {
	return m.repos, m.err
}
func (m *mockStore) Get(id string) (*pkg.Repository, error) { return m.repo, m.err }

func setupController(store repository.Store, m func(http.Handler) http.Handler) *bastion.Bastion {
	updatedDESC := sorting.NewSort("updated_at_desc", "updated_at DESC", "Updated date descending")
	updatedASC := sorting.NewSort("updated_at_asc", "updated_at ASC", "Updated date ascendant")
	createdDESC := sorting.NewSort("created_at_desc", "created_at DESC", "Created date descending")
	createdASC := sorting.NewSort("created_at_asc", "created_at ASC", "Created date ascendant")

	publicVisibility := filtering.NewValue("public", "public repos")
	privateVisibility := filtering.NewValue("private", "private repos")
	visibilityFilter := filtering.NewText("visibility", "filters the repos by their visibility", publicVisibility, privateVisibility)

	listing := middleware.Listing(
		middleware.MaxAllowedLimit(50),
		middleware.Sort(updatedDESC, updatedASC, createdDESC, createdASC),
		middleware.Filter(visibilityFilter),
	)

	s := repository.Service{Store: store}
	app := bastion.New()
	app.APIRouter.Use(m)
	app.APIRouter.Use(listing)
	app.APIRouter.Get("/user/repos/", repository.ListingOwnRepos(s))

	return app
}

func TestListOwnerReposWhenEmpty(t *testing.T) {
	t.Parallel()
	store := &mockStore{repos: []pkg.Repository{}}
	app := setupController(store, authRequest)

	e := bastion.Tester(t, app)
	res := e.GET("/user/repos/").Expect().
		JSON().Object()
	res.Value("results").Array().Empty()
	res.Value("listing").Object().
		ContainsKey("paging").
		ContainsKey("sorting").
		ContainsKey("filtering")
}

func TestListOwnerReposWithValues(t *testing.T) {
	t.Parallel()
	store := &mockStore{repos: []pkg.Repository{
		{Name: "test public", Visibility: pkg.Public},
		{Name: "test private", Visibility: pkg.Private},
	}}
	app := setupController(store, authRequest)

	e := bastion.Tester(t, app)
	res := e.GET("/user/repos/").
		Expect().
		Status(http.StatusOK).JSON().Object()

	results := res.Value("results").Array().NotEmpty()
	results.Length().Equal(2)
	results.First().Object().
		ContainsKey("id").
		ContainsKey("name").
		ContainsKey("current_branch").
		ContainsKey("visibility").
		ContainsKey("createdAt").
		ContainsKey("updatedAt").
		ContainsKey("owner")
}

// TODO: this should be an integral test
//func TestListOwnerReposWithValuesFilter(t *testing.T) {
//	t.Parallel()
//
//	tt := []struct {
//		name                string
//		params              string
//		amount              int
//		repoExpectedResults map[string]interface{}
//	}{
//		{
//			"filter public repos",
//			"visibility=public",
//			1,
//			map[string]interface{}{
//				"name":           "test public",
//				"visibility":     "public",
//				"current_branch": "master",
//				"owner":          tempUser.ID,
//			},
//		},
//		{
//			"filter private repos",
//			"visibility=private",
//			1,
//			map[string]interface{}{
//				"name":           "test private",
//				"visibility":     "private",
//				"current_branch": "master",
//				"owner":          tempUser.ID,
//			},
//		},
//	}
//
//	store := &mockStore{repos: []pkg.Repository{
//		{Name: "test public", Visibility: pkg.Public},
//		{Name: "test private", Visibility: pkg.Private},
//	}}
//	app := setupController(store, authRequest)
//	e := bastion.Tester(t, app)
//	for _, tc := range tt {
//		t.Run(tc.name, func(t *testing.T) {
//			res := e.GET("/user/repos/").WithQueryString(tc.params).
//				Expect().
//				Status(http.StatusOK).JSON().Object()
//
//			results := res.Value("results").Array()
//			results.Length().Equal(tc.amount)
//			results.First().Object().
//				ValueEqual("name", tc.repoExpectedResults["name"]).
//				ValueEqual("current_branch", tc.repoExpectedResults["current_branch"]).
//				ValueEqual("visibility", tc.repoExpectedResults["visibility"]).
//				ValueEqual("owner", tc.repoExpectedResults["owner"])
//		})
//	}
//}
