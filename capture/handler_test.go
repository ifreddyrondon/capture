package capture_test

import (
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/araddon/dateparse"
	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/render/json"
	"github.com/ifreddyrondon/gocapture/app"
	"github.com/ifreddyrondon/gocapture/capture"
	"github.com/ifreddyrondon/gocapture/database"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

var once sync.Once
var db *gorm.DB

func getDB() *gorm.DB {
	once.Do(func() {
		ds := database.Open("postgres://localhost/captures_app_test?sslmode=disable")
		db = ds.DB
	})
	return db
}

func setup(t *testing.T) (*bastion.Bastion, func()) {
	t.Parallel()

	// get a random table to allow parallel execution
	schemaName := fmt.Sprintf("captures_%v", time.Now().UnixNano())
	service := capture.PGService{DB: getDB().Table(schemaName)}
	service.Migrate()
	teardown := func() { service.Drop() }

	handler := capture.Handler{
		Service: &service,
		Render:  json.NewRender,
		CtxKey:  app.ContextKey("capture"),
	}

	app := bastion.New(bastion.Options{})
	app.APIRouter.Mount("/captures/", handler.Router())

	return app, teardown
}

func TestCreateValidCapture(t *testing.T) {
	app, teardown := setup(t)
	defer teardown()

	e := bastion.Tester(t, app)
	tt := []struct {
		name     string
		payload  map[string]interface{}
		response map[string]interface{}
	}{
		{
			name: "create capture with payload, date and point",
			payload: map[string]interface{}{
				"latitude":  1,
				"longitude": 12,
				"timestamp": "630655260",
				"payload":   map[string]interface{}{"power": []interface{}{-70.0, -100.1, 3.1}},
			},
			response: map[string]interface{}{
				"lat":       1.0,
				"lng":       12.0,
				"timestamp": "1989-12-26T06:01:00Z",
				"payload":   map[string]interface{}{"power": []interface{}{-70.0, -100.1, 3.1}},
			},
		},
		{
			name: "create capture with payload and date without point",
			payload: map[string]interface{}{
				"date":    "630655260",
				"payload": map[string]interface{}{"power": []interface{}{-70.0, -100.1, 3.1}},
			},
			response: map[string]interface{}{
				"lat":       nil,
				"lng":       nil,
				"timestamp": "1989-12-26T06:01:00Z",
				"payload":   map[string]interface{}{"power": []interface{}{-70.0, -100.1, 3.1}},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			e.POST("/captures/").
				WithJSON(tc.payload).
				Expect().
				Status(http.StatusCreated).
				JSON().Object().
				ContainsKey("payload").ValueEqual("payload", tc.response["payload"]).
				ContainsKey("lat").ValueEqual("lat", tc.response["lat"]).
				ContainsKey("lng").ValueEqual("lng", tc.response["lng"]).
				ContainsKey("timestamp").ValueEqual("timestamp", tc.response["timestamp"]).
				ContainsKey("id").NotEmpty().
				ContainsKey("createdAt").NotEmpty().
				ContainsKey("updatedAt").NotEmpty()
		})
	}
}

func TestCreateOnlyPayloadCapture(t *testing.T) {
	app, teardown := setup(t)
	defer teardown()

	payload := map[string]interface{}{
		"payload": map[string]interface{}{"power": []interface{}{-70.0, -100.1, 3.1}},
	}
	response := map[string]interface{}{
		"payload": map[string]interface{}{"power": []interface{}{-70.0, -100.1, 3.1}},
	}
	e := bastion.Tester(t, app)
	e.POST("/captures/").
		WithJSON(payload).
		Expect().
		Status(http.StatusCreated).
		JSON().Object().
		ContainsKey("payload").ValueEqual("payload", response["payload"]).
		ContainsKey("lat").ValueEqual("lat", nil).
		ContainsKey("lng").ValueEqual("lng", nil).
		ContainsKey("timestamp").NotEmpty().
		ContainsKey("id").NotEmpty().
		ContainsKey("createdAt").NotEmpty().
		ContainsKey("updatedAt").NotEmpty()
}

func TestCreateInvalidCapture(t *testing.T) {
	app, teardown := setup(t)
	defer teardown()

	e := bastion.Tester(t, app)
	tt := []struct {
		name     string
		payload  map[string]interface{}
		response map[string]interface{}
	}{
		{
			name:    "bad request, missing body",
			payload: map[string]interface{}{},
			response: map[string]interface{}{
				"status":  400.0,
				"error":   "Bad Request",
				"message": "missing payload value",
			},
		},
		{
			name:    "bad request, missing lng",
			payload: map[string]interface{}{"lat": 1, "date": "630655260"},
			response: map[string]interface{}{
				"status":  400.0,
				"error":   "Bad Request",
				"message": "missing longitude",
			},
		},
		{
			name:    "bad request, missing lat",
			payload: map[string]interface{}{"lng": 1, "date": "630655260"},
			response: map[string]interface{}{
				"status":  400.0,
				"error":   "Bad Request",
				"message": "missing latitude",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			e.POST("/captures/").
				WithJSON(tc.payload).
				Expect().
				Status(http.StatusBadRequest).
				JSON().Object().Equal(tc.response)
		})
	}
}

func TestCreateInvalidPayloadCapture(t *testing.T) {
	response := map[string]interface{}{
		"status":  400.0,
		"error":   "Bad Request",
		"message": "cannot unmarshal json into valid capture",
	}

	app, teardown := setup(t)
	defer teardown()
	e := bastion.Tester(t, app)
	e.POST("/captures/").
		WithJSON("{").
		Expect().
		Status(http.StatusBadRequest).
		JSON().Object().Equal(response)
}

func TestListCapturesWhenEmpty(t *testing.T) {
	app, teardown := setup(t)
	defer teardown()

	e := bastion.Tester(t, app)
	e.GET("/captures/").Expect().JSON().Array().Empty()
}

func TestListCapturesWithValues(t *testing.T) {
	app, teardown := setup(t)
	defer teardown()

	payload := map[string]interface{}{
		"payload": map[string]interface{}{"power": []interface{}{-70.0, -100.1, 3.1}},
	}
	e := bastion.Tester(t, app)
	e.POST("/captures/").WithJSON(payload).Expect().Status(http.StatusCreated)

	array := e.GET("/captures/").
		Expect().
		Status(http.StatusOK).
		JSON().Array().NotEmpty()

	array.Length().Equal(1)
	array.First().Object().
		ContainsKey("payload").
		ContainsKey("lat").
		ContainsKey("lng").
		ContainsKey("timestamp").
		ContainsKey("id").
		ContainsKey("createdAt").
		ContainsKey("updatedAt")
}

func TestGetCapture(t *testing.T) {
	app, teardown := setup(t)
	defer teardown()

	capPayload := map[string]interface{}{
		"lat":       1.0,
		"lng":       12.0,
		"timestamp": "1989-12-26T06:01:00Z",
		"payload":   map[string]interface{}{"power": []interface{}{-70.0, -100.1, 3.1}},
	}
	e := bastion.Tester(t, app)
	obj := e.POST("/captures/").WithJSON(capPayload).Expect().Status(http.StatusCreated).
		JSON().Object().Raw()

	e.GET(fmt.Sprintf("/captures/%v", obj["id"])).Expect().
		Status(http.StatusOK).
		JSON().Object().
		ContainsKey("payload").ValueEqual("payload", capPayload["payload"]).
		ContainsKey("lat").ValueEqual("lat", capPayload["lat"]).
		ContainsKey("lng").ValueEqual("lng", capPayload["lng"]).
		ContainsKey("timestamp").NotEmpty().
		ContainsKey("id").NotEmpty().
		ContainsKey("createdAt").NotEmpty().
		ContainsKey("updatedAt").NotEmpty()
}

func TestGetMissingCapture(t *testing.T) {
	app, teardown := setup(t)
	defer teardown()

	response := map[string]interface{}{
		"status":  404.0,
		"error":   "Not Found",
		"message": "not found capture",
	}

	e := bastion.Tester(t, app)
	e.GET(fmt.Sprint("/captures/123123")).Expect().
		Status(http.StatusNotFound).
		JSON().Object().Equal(response)
}

func TestGetCaptureBadRequest(t *testing.T) {
	app, teardown := setup(t)
	defer teardown()

	response := map[string]interface{}{
		"status":  400.0,
		"error":   "Bad Request",
		"message": "invalid capture id",
	}

	e := bastion.Tester(t, app)
	e.GET(fmt.Sprint("/captures/ads")).Expect().
		Status(http.StatusBadRequest).
		JSON().Object().Equal(response)
}

func TestDeleteCapture(t *testing.T) {
	app, teardown := setup(t)
	defer teardown()

	capPayload := map[string]interface{}{
		"lat":       1.0,
		"lng":       12.0,
		"timestamp": "1989-12-26T06:01:00Z",
		"payload":   map[string]interface{}{"power": []interface{}{-70.0, -100.1, 3.1}},
	}
	e := bastion.Tester(t, app)
	obj := e.POST("/captures/").WithJSON(capPayload).Expect().Status(http.StatusCreated).
		JSON().Object().Raw()

	id := obj["id"]

	e.GET(fmt.Sprintf("/captures/%v", id)).Expect().Status(http.StatusOK)
	e.DELETE(fmt.Sprintf("/captures/%v", id)).Expect().Status(http.StatusNoContent)
	e.GET(fmt.Sprintf("/captures/%v", id)).Expect().Status(http.StatusNotFound)
}

func TestUpdateCapture(t *testing.T) {
	app, teardown := setup(t)
	defer teardown()

	e := bastion.Tester(t, app)
	capPayload := map[string]interface{}{
		"lat":       1.0,
		"lng":       12.0,
		"timestamp": "1989-12-26T06:01:00Z",
		"payload":   map[string]interface{}{"power": []interface{}{-70.0, -100.1, 3.1}},
	}

	tt := []struct {
		name          string
		updatePayload map[string]interface{}
	}{
		{
			"update lat",
			map[string]interface{}{
				"lat":       89.0,
				"lng":       capPayload["lng"],
				"timestamp": capPayload["timestamp"],
				"payload":   capPayload["payload"],
			},
		},
		{
			"update lng",
			map[string]interface{}{
				"lat":       capPayload["lat"],
				"lng":       30,
				"timestamp": capPayload["timestamp"],
				"payload":   capPayload["payload"],
			},
		},
		{
			"update timestamp",
			map[string]interface{}{
				"lat":       capPayload["lat"],
				"lng":       capPayload["lng"],
				"timestamp": "2006-07-12T06:01:00Z",
				"payload":   capPayload["payload"],
			},
		},
		{
			"update payload",
			map[string]interface{}{
				"lat":       capPayload["lat"],
				"lng":       capPayload["lng"],
				"timestamp": capPayload["timestamp"],
				"payload":   map[string]interface{}{"power": []interface{}{1}},
			},
		},
		{
			"do not update id",
			map[string]interface{}{
				"id":        "123",
				"lat":       capPayload["lat"],
				"lng":       capPayload["lng"],
				"timestamp": capPayload["timestamp"],
				"payload":   capPayload["payload"],
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			createdObj := e.POST("/captures/").WithJSON(capPayload).Expect().
				Status(http.StatusCreated).JSON().Object().Raw()
			e.GET(fmt.Sprintf("/captures/%v", createdObj["id"])).Expect().Status(http.StatusOK)
			tc.updatePayload["id"] = createdObj["id"]
			updatedObj := e.PUT(fmt.Sprintf("/captures/%v", createdObj["id"])).WithJSON(tc.updatePayload).Expect().
				Status(http.StatusOK).
				JSON().Object().
				ContainsKey("id").ValueEqual("id", tc.updatePayload["id"]).
				ContainsKey("lat").ValueEqual("lat", tc.updatePayload["lat"]).
				ContainsKey("lng").ValueEqual("lng", tc.updatePayload["lng"]).
				ContainsKey("timestamp").ValueEqual("timestamp", tc.updatePayload["timestamp"]).
				ContainsKey("payload").ValueEqual("payload", tc.updatePayload["payload"]).
				ContainsKey("createdAt").NotEmpty().
				ContainsKey("updatedAt").NotEmpty().
				Raw()

			// TODO: createdAt should be equal

			updatedAtFromCreate, err := dateparse.ParseAny(createdObj["updatedAt"].(string))
			assert.Nil(t, err)
			updatedAtFromUpdate, err := dateparse.ParseAny(updatedObj["updatedAt"].(string))
			assert.Nil(t, err)
			updatedAtFromUpdate.After(updatedAtFromCreate)
		})
	}
}

func TestUpdateCaptureFailsBadRequest(t *testing.T) {
	app, teardown := setup(t)
	defer teardown()

	e := bastion.Tester(t, app)
	okPayload := map[string]interface{}{
		"lat":       1.0,
		"lng":       12.0,
		"timestamp": "1989-12-26T06:01:00Z",
		"payload":   map[string]interface{}{"power": []interface{}{-70.0, -100.1, 3.1}},
	}

	tt := []struct {
		name          string
		updatePayload map[string]interface{}
		response      map[string]interface{}
	}{
		{
			"lat out of range",
			map[string]interface{}{
				"lat":       200.0,
				"lng":       12.0,
				"timestamp": "1989-12-26T06:01:00Z",
				"payload":   map[string]interface{}{"power": []interface{}{-70.0, -100.1, 3.1}},
			},
			map[string]interface{}{
				"status":  400.0,
				"error":   "Bad Request",
				"message": "latitude out of boundaries, may range from -90.0 to 90.0",
			},
		},
		{
			"missing payload",
			map[string]interface{}{
				"lat":       1.0,
				"lng":       12.0,
				"timestamp": "1989-12-26T06:01:00Z",
			},
			map[string]interface{}{
				"status":  400.0,
				"error":   "Bad Request",
				"message": "missing payload value",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			createdObj := e.POST("/captures/").WithJSON(okPayload).Expect().
				Status(http.StatusCreated).JSON().Object().Raw()
			e.GET(fmt.Sprintf("/captures/%v", createdObj["id"])).Expect().Status(http.StatusOK)
			tc.updatePayload["id"] = createdObj["id"]
			e.PUT(fmt.Sprintf("/captures/%v", createdObj["id"])).WithJSON(tc.updatePayload).Expect().
				Status(http.StatusBadRequest).
				JSON().Object().Equal(tc.response)
		})
	}
}

func TestUnmarshalCapturesFail(t *testing.T) {
	app, teardown := setup(t)
	defer teardown()

	e := bastion.Tester(t, app)
	tt := []struct {
		name     string
		payload  []interface{}
		response map[string]interface{}
	}{
		{
			name:    "bad request, missing body",
			payload: []interface{}{},
			response: map[string]interface{}{
				"status":  400.0,
				"error":   "Bad Request",
				"message": "cannot unmarshal json into valid captures, it needs at least one capture",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			e.POST("/captures/").
				WithJSON(tc.payload).
				Expect().
				Status(http.StatusBadRequest).
				JSON().Object().Equal(tc.response)
		})
	}
}
