package rest_test

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/capture/pkg/creating"
	"github.com/ifreddyrondon/capture/pkg/domain"
	"github.com/ifreddyrondon/capture/pkg/http/rest"
)

type mockCreatingService struct {
	repo *creating.Repository
	err  error
}

func (m *mockCreatingService) CreateRepo(*domain.User, creating.Payload) (*creating.Repository, error) {
	return m.repo, m.err
}

func setupCreateHandler(s creating.Service, m func(http.Handler) http.Handler) *bastion.Bastion {
	app := bastion.New()
	app.APIRouter.Use(m)
	app.APIRouter.Post("/", rest.Creating(s))
	return app
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
	app := setupCreateHandler(s, withUserMiddle(defaultUser))
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
	app := setupCreateHandler(s, withUserMiddle(defaultUser))
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

func TestCreateRepoFailInternalErrorGettingUser(t *testing.T) {
	t.Parallel()
	s := &mockCreatingService{}
	app := setupCreateHandler(s, withUserMiddle(nil))
	e := bastion.Tester(t, app)

	payload := map[string]interface{}{"name": "test"}
	response := map[string]interface{}{
		"status":  500.0,
		"error":   "Internal Server Error",
		"message": "looks like something went wrong",
	}

	e.POST("/").
		WithJSON(payload).
		Expect().
		Status(http.StatusInternalServerError).
		JSON().Object().Equal(response)
}

func TestCreateRepoFailInternalErrorCreatingRepo(t *testing.T) {
	t.Parallel()
	s := &mockCreatingService{err: errors.New("test")}
	app := setupCreateHandler(s, withUserMiddle(defaultUser))
	e := bastion.Tester(t, app)

	payload := map[string]interface{}{"name": "test"}
	response := map[string]interface{}{
		"status":  500.0,
		"error":   "Internal Server Error",
		"message": "looks like something went wrong",
	}

	e.POST("/").
		WithJSON(payload).
		Expect().
		Status(http.StatusInternalServerError).
		JSON().Object().Equal(response)
}
