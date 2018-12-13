package rest_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/render"
	"github.com/ifreddyrondon/capture/pkg"
	"github.com/ifreddyrondon/capture/pkg/creating"
	"github.com/ifreddyrondon/capture/pkg/http/rest"
	"github.com/ifreddyrondon/capture/pkg/http/rest/middleware"
)

var tempUser = pkg.User{Email: "test@example.com", ID: "0162eb39-a65e-04a1-7ad9-d663bb49a396"}

func notAuthRequest(next http.Handler) http.Handler {
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

func authRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := middleware.WithUser(r.Context(), &tempUser)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type mockCreatingService struct {
	repo *creating.Repository
	err  error
}

func (m *mockCreatingService) CreateRepo(*pkg.User, creating.Payload) (*creating.Repository, error) {
	return m.repo, m.err
}

func TestCreateRepoSuccess(t *testing.T) {
	t.Parallel()

	repo := &creating.Repository{
		ID:         "01679604-d8f6-29ce-2fe2-5d66dfa2d194",
		Name:       "test",
		Visibility: "public",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	s := &mockCreatingService{repo: repo}
	app := bastion.New()
	app.APIRouter.Mount("/", rest.Creating(s, authRequest))
	e := bastion.Tester(t, app)

	payload := map[string]interface{}{"name": "test"}

	e.POST("/").
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

func TestCreateRepoFailBadRequest(t *testing.T) {
	t.Parallel()
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

	s := &mockCreatingService{}
	app := bastion.New()
	app.APIRouter.Mount("/", rest.Creating(s, authRequest))
	e := bastion.Tester(t, app)
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			e.POST("/").
				WithJSON(tc.payload).
				Expect().
				Status(http.StatusBadRequest).
				JSON().Object().Equal(tc.response)
		})
	}
}
