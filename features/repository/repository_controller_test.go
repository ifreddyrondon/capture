package repository_test

import (
	"errors"
	"net/http"
	"testing"

	"bytes"

	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/render"
	"github.com/ifreddyrondon/capture/config"
	"github.com/ifreddyrondon/capture/features"
	"github.com/ifreddyrondon/capture/features/repository"
	"github.com/ifreddyrondon/capture/features/user"
	"github.com/jinzhu/gorm"
	"gopkg.in/src-d/go-kallax.v1"
)

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
		ctx := user.WithUser(r.Context(), &features.User{Email: "test@example.com", ID: kallax.NewULID()})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type mockStore struct {
	repos []features.Repository
	err   error
}

func (m *mockStore) Save(u *features.User, c *features.Repository) error { return m.err }
func (m *mockStore) List(l repository.ListingRepo) ([]features.Repository, error) {
	return m.repos, m.err
}

func setup(t *testing.T, isAuth func(http.Handler) http.Handler) (*bastion.Bastion, func()) {
	toml := []byte(`PG="postgres://localhost/captures_app_test?sslmode=disable"`)
	cfg, err := config.New(config.Source(bytes.NewBuffer(toml)))
	if err != nil {
		t.Fatal(err)
	}

	db := cfg.Resources.Get("database").(*gorm.DB)
	store := repository.NewPGStore(db.Table("repositories"))
	store.Migrate()

	app := bastion.New()
	app.APIRouter.Mount("/repository/", repository.Routes(store, isAuth, loggedUser))

	return app, func() { store.Drop() }
}

func TestCreateRepositorySuccess(t *testing.T) {
	app, teardown := setup(t, authOK)
	defer teardown()

	e := bastion.Tester(t, app)
	payload := map[string]interface{}{"name": "test"}

	e.POST("/repository/").
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
	app, teardown := setup(t, authOK)
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
			e.POST("/repository/").
				WithJSON(tc.payload).
				Expect().
				Status(http.StatusBadRequest).
				JSON().Object().Equal(tc.response)
		})
	}
}

func setupController(store repository.Store, isAuth func(http.Handler) http.Handler) *bastion.Bastion {
	app := bastion.New()
	app.APIRouter.Mount("/repository/", repository.Routes(store, isAuth, loggedUser))

	return app
}

func TestCreateRepositorySaveFail(t *testing.T) {
	t.Parallel()
	store := &mockStore{err: errors.New("test")}
	app := setupController(store, authOK)

	e := bastion.Tester(t, app)
	payload := map[string]interface{}{"name": "test"}

	e.POST("/repository/").
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

	e.POST("/repository/").
		WithJSON(payload).
		Expect().
		Status(http.StatusForbidden).
		JSON().Object().Equal(response)
}
