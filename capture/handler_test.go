package capture_test

import (
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	mgo "gopkg.in/mgo.v2"

	"fmt"

	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/render/json"
	"github.com/ifreddyrondon/gocapture/app"
	"github.com/ifreddyrondon/gocapture/capture"
	"github.com/ifreddyrondon/gocapture/database"
)

var once sync.Once
var db *mgo.Database

func getDB(t *testing.T) *mgo.Database {
	once.Do(func() {
		ds, err := database.Open("localhost/captures_test")
		db = ds.DB()
		require.Nil(t, err)
	})
	return db
}

func setup(t *testing.T) (*bastion.Bastion, func()) {
	t.Parallel()

	// get a random collection to allow parallel execution
	collection := getDB(t).C(fmt.Sprintf("captures.%v", time.Now().UnixNano()))
	teardown := func() { collection.DropCollection() }

	service := capture.MgoService{Collection: collection}
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

	tt := []struct {
		name     string
		payload  map[string]interface{}
		response map[string]interface{}
	}{
		{
			"create capture with date name",
			map[string]interface{}{"lat": 1, "lng": 12, "date": "1989-12-26T06:01:00.00Z"},
			map[string]interface{}{
				"payload":   nil,
				"lat":       1.0,
				"lng":       12.0,
				"timestamp": "1989-12-26T06:01:00Z",
			},
		},
		{
			name:    "create capture with timestamp name",
			payload: map[string]interface{}{"lat": 1, "lng": 12, "timestamp": "630655260"},
			response: map[string]interface{}{
				"payload":   nil,
				"lat":       1.0,
				"lng":       12.0,
				"timestamp": "1989-12-26T06:01:00Z",
			},
		},
		{
			name:    "create capture with latitude, longitude and data names",
			payload: map[string]interface{}{"latitude": 1, "longitude": 12, "date": "630655260"},
			response: map[string]interface{}{
				"payload":   nil,
				"lat":       1.0,
				"lng":       12.0,
				"timestamp": "1989-12-26T06:01:00Z",
			},
		},
		{
			name: "create capture with payload",
			payload: map[string]interface{}{
				"latitude":  1,
				"longitude": 12,
				"date":      "630655260",
				"payload":   []float32{-78.75, -80.5, -73.75, -70.75, -72},
			},
			response: map[string]interface{}{
				"lat":       1.0,
				"lng":       12.0,
				"timestamp": "1989-12-26T06:01:00Z",
				"payload":   []float32{-78.75, -80.5, -73.75, -70.75, -72},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			e := bastion.Tester(t, app)
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
				ContainsKey("created_date").NotEmpty().
				ContainsKey("last_modified").NotEmpty()
		})
	}
}

func TestCreateInValidCapture(t *testing.T) {
	app, teardown := setup(t)
	defer teardown()

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
				"message": "missing latitude",
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
			e := bastion.Tester(t, app)
			e.POST("/captures/").
				WithJSON(tc.payload).
				Expect().
				Status(http.StatusBadRequest).
				JSON().Object().Equal(tc.response)
		})
	}
}

func TestCreateInValidPayloadCapture(t *testing.T) {
	response := map[string]interface{}{
		"status":  400.0,
		"error":   "Bad Request",
		"message": "cannot unmarshal json into Point value",
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

	payload := map[string]interface{}{"lat": 1, "lng": 12, "date": "1989-12-26T06:01:00.00Z"}
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
		ContainsKey("created_date").
		ContainsKey("last_modified")
}

func TestGetCapture(t *testing.T) {
	app, teardown := setup(t)
	defer teardown()

	capPayload := map[string]interface{}{
		"lat":       1.0,
		"lng":       12.0,
		"timestamp": "1989-12-26T06:01:00Z",
		"payload":   []float32{-78.75, -80.5, -73.75, -70.75, -72},
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
		ContainsKey("created_date").NotEmpty().
		ContainsKey("last_modified").NotEmpty()
}

func TestGetMissingCapture(t *testing.T) {
	app, teardown := setup(t)
	defer teardown()

	response := map[string]interface{}{
		"status":  404.0,
		"error":   "Not Found",
		"message": "not found",
	}

	e := bastion.Tester(t, app)
	e.GET(fmt.Sprint("/captures/5ab3a603841d0925708a6ea7")).Expect().
		Status(http.StatusNotFound).
		JSON().Object().Equal(response)
}
