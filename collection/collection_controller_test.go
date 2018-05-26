package collection_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/render/json"
	"github.com/ifreddyrondon/gocapture/collection"
)

type MockService struct{}

func (r *MockService) Save(c *collection.Collection) error {
	return errors.New("test")
}

func setupControllerMockService() *bastion.Bastion {
	service := &MockService{}

	controller := collection.NewController(service, json.NewRender)

	app := bastion.New(bastion.Options{})
	app.APIRouter.Mount("/collection/", controller.Router())
	return app
}

func setupController(t *testing.T) (*bastion.Bastion, func()) {
	service, teardown := setupService(t)

	controller := collection.NewController(service, json.NewRender)

	app := bastion.New(bastion.Options{})
	app.APIRouter.Mount("/collection/", controller.Router())

	return app, teardown
}

func TestCreateCollectionSuccess(t *testing.T) {
	app, teardown := setupController(t)
	defer teardown()

	e := bastion.Tester(t, app)
	payload := map[string]interface{}{"name": "test"}

	e.POST("/collection/").
		WithJSON(payload).
		Expect().
		Status(http.StatusCreated).
		JSON().Object().
		ContainsKey("name").ValueEqual("name", payload["name"]).
		ContainsKey("shared").ValueEqual("shared", false).
		ContainsKey("id").NotEmpty().
		ContainsKey("createdAt").NotEmpty().
		ContainsKey("updatedAt").NotEmpty()
}

func TestCreateCollectionFail(t *testing.T) {
	app, teardown := setupController(t)
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
			e.POST("/collection/").
				WithJSON(tc.payload).
				Expect().
				Status(http.StatusBadRequest).
				JSON().Object().Equal(tc.response)
		})
	}
}

func TestCreateCollectionSaveFail(t *testing.T) {
	t.Parallel()

	app := setupControllerMockService()

	e := bastion.Tester(t, app)
	payload := map[string]interface{}{"name": "test"}

	e.POST("/collection/").
		WithJSON(payload).
		Expect().
		Status(http.StatusInternalServerError).
		JSON().Object()
}
