package repository_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/render"
	"github.com/ifreddyrondon/capture/config"
	"github.com/ifreddyrondon/capture/features"
	"github.com/ifreddyrondon/capture/features/repository"
	"github.com/ifreddyrondon/capture/features/user"
	"gopkg.in/src-d/go-kallax.v1"
)

var tempUser = features.User{Email: "test@example.com", ID: kallax.NewULID()}

func authOK(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func authFails(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := render.HTTPError{
			Status:  http.StatusForbidden,
			Error:   http.StatusText(http.StatusForbidden),
			Message: "you don’t have permission to access this resource",
		}
		render.NewJSON().Response(w, http.StatusForbidden, err)
		return
	})
}

func loggedUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := user.WithUser(r.Context(), &tempUser)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type mockStore struct {
	repo  *features.Repository
	repos []features.Repository
	err   error
}

func (m *mockStore) Drop()                                               {}
func (m *mockStore) Save(u *features.User, c *features.Repository) error { return m.err }
func (m *mockStore) List(l repository.ListingRepo) ([]features.Repository, error) {
	return m.repos, m.err
}
func (m *mockStore) Get(id kallax.ULID) (*features.Repository, error) { return m.repo, m.err }

func setupUserController(t *testing.T, isAuth func(http.Handler) http.Handler) (*bastion.Bastion, func()) {
	cfg, err := config.FromString(`PG="postgres://localhost/captures_app_test?sslmode=disable"`)
	if err != nil {
		t.Fatal(err)
	}

	store := cfg.Resources.Get("repository-store").(repository.Store)
	service := repository.Service{Store: store}
	app := bastion.New()
	app.APIRouter.Mount("/user/repos/", repository.UserRoutes(service, isAuth, loggedUser))

	return app, func() { store.Drop() }
}

func TestCreateRepositorySuccess(t *testing.T) {
	app, teardown := setupUserController(t, authOK)
	defer teardown()

	e := bastion.Tester(t, app)
	payload := map[string]interface{}{"name": "test"}

	e.POST("/user/repos/").
		WithJSON(payload).
		Expect().
		Status(http.StatusCreated).
		JSON().Object().
		ContainsKey("name").ValueEqual("name", payload["name"]).
		ContainsKey("visibility").ValueEqual("visibility", "public").
		ContainsKey("id").NotEmpty().
		ContainsKey("createdAt").NotEmpty().
		ContainsKey("updatedAt").NotEmpty()
}

func TestCreateRepositoryFail(t *testing.T) {
	app, teardown := setupUserController(t, authOK)
	defer teardown()

	e := bastion.Tester(t, app)
	tt := []struct {
		name     string
		payload  map[string]interface{}
		response map[string]interface{}
	}{
		{
			name:    "no data",
			payload: map[string]interface{}{},
			response: map[string]interface{}{
				"status":  400.0,
				"error":   "Bad Request",
				"message": "name must not be blank",
			},
		},
		{
			name:    "empty name",
			payload: map[string]interface{}{"name": ""},
			response: map[string]interface{}{
				"status":  400.0,
				"error":   "Bad Request",
				"message": "name must not be blank",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			e.POST("/user/repos/").
				WithJSON(tc.payload).
				Expect().
				Status(http.StatusBadRequest).
				JSON().Object().Equal(tc.response)
		})
	}
}

func setupController(store repository.Store, isAuth func(http.Handler) http.Handler) *bastion.Bastion {
	service := repository.Service{Store: store}
	app := bastion.New()
	app.APIRouter.Mount("/user/repos/", repository.UserRoutes(service, isAuth, loggedUser))

	return app
}

func TestCreateRepositorySaveFail(t *testing.T) {
	t.Parallel()
	store := &mockStore{err: errors.New("test")}
	app := setupController(store, authOK)

	e := bastion.Tester(t, app)
	payload := map[string]interface{}{"name": "test"}

	e.POST("/user/repos/").
		WithJSON(payload).
		Expect().
		Status(http.StatusInternalServerError).
		JSON().Object()
}

func TestCreateRepositoryNotAuthorized(t *testing.T) {
	t.Parallel()
	app := setupController(&mockStore{}, authFails)

	response := map[string]interface{}{
		"status":  403.0,
		"error":   "Forbidden",
		"message": "you don’t have permission to access this resource",
	}

	e := bastion.Tester(t, app)
	payload := map[string]interface{}{"name": "test"}

	e.POST("/user/repos/").
		WithJSON(payload).
		Expect().
		Status(http.StatusForbidden).
		JSON().Object().Equal(response)
}

func TestListOwnerReposWhenEmpty(t *testing.T) {
	t.Parallel()
	store := &mockStore{repos: []features.Repository{}}
	app := setupController(store, authOK)

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
	app, teardown := setupUserController(t, authOK)
	defer teardown()

	public := map[string]interface{}{"name": "test public"}
	private := map[string]interface{}{"name": "test private", "visibility": "private"}

	e := bastion.Tester(t, app)
	e.POST("/user/repos/").WithJSON(public).Expect().Status(http.StatusCreated)
	e.POST("/user/repos/").WithJSON(private).Expect().Status(http.StatusCreated)

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

func TestListOwnerReposWithValuesFilter(t *testing.T) {
	app, teardown := setupUserController(t, authOK)
	defer teardown()

	public := map[string]interface{}{"name": "test public"}
	private := map[string]interface{}{"name": "test private", "visibility": "private"}

	e := bastion.Tester(t, app)
	e.POST("/user/repos/").WithJSON(public).Expect().Status(http.StatusCreated)
	e.POST("/user/repos/").WithJSON(private).Expect().Status(http.StatusCreated)

	tt := []struct {
		name                string
		params              string
		amount              int
		repoExpectedResults map[string]interface{}
	}{
		{
			"filter public repos",
			"visibility=public",
			1,
			map[string]interface{}{
				"name":           "test public",
				"visibility":     "public",
				"current_branch": "master",
				"owner":          tempUser.ID.String(),
			},
		},
		{
			"filter private repos",
			"visibility=private",
			1,
			map[string]interface{}{
				"name":           "test private",
				"visibility":     "private",
				"current_branch": "master",
				"owner":          tempUser.ID.String(),
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			res := e.GET("/user/repos/").WithQueryString(tc.params).
				Expect().
				Status(http.StatusOK).JSON().Object()

			results := res.Value("results").Array()
			results.Length().Equal(tc.amount)
			results.First().Object().
				ValueEqual("name", tc.repoExpectedResults["name"]).
				ValueEqual("current_branch", tc.repoExpectedResults["current_branch"]).
				ValueEqual("visibility", tc.repoExpectedResults["visibility"]).
				ValueEqual("owner", tc.repoExpectedResults["owner"])
		})
	}
}
