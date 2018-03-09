package capture_test

import (
	"net/http"
	"testing"

	"fmt"

	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/render/json"
	"github.com/ifreddyrondon/gocapture/capture"
	"github.com/ifreddyrondon/gocapture/database"
)

func setup(t *testing.T) (*bastion.Bastion, func()) {
	t.Parallel()

	ds, err := database.Open("localhost/captures_test")
	if err != nil {
		t.Fatalf("could not create database, err: %v", err)
	}
	db := ds.DB()
	service := capture.MgoService{DB: db}
	handler := capture.Handler{
		Service: &service,
		Render:  json.NewRender,
	}

	app := bastion.New(nil)
	app.APIRouter.Mount(fmt.Sprintf("/%v/", handler.Pattern()), handler.Router())

	teardown := func() { ds.DB().DropDatabase() }

	return app, teardown
}

func TestCreateValidCapture(t *testing.T) {
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
			app, teardown := setup(t)
			defer teardown()

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
			app, teardown := setup(t)
			defer teardown()

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
