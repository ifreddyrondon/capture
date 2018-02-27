package capture_test

import (
	"log"
	"net/http"
	"testing"

	"os"

	"time"

	"github.com/ifreddyrondon/gobastion"
	"github.com/ifreddyrondon/gocapture/app"
	"github.com/ifreddyrondon/gocapture/capture"
	"github.com/ifreddyrondon/gocapture/database"
	"gopkg.in/mgo.v2"
)

var bastion *gobastion.Bastion
var db *mgo.Database

func clearCollection() {
	db.DropDatabase()
}

func TestMain(m *testing.M) {
	reader := new(gobastion.JsonReader)
	responder := new(gobastion.JsonResponder)

	ds, err := database.Open("localhost/captures_test")
	if err != nil {
		log.Panic(err)
	}

	service := capture.MgoService{DB: ds.DB()}
	handler := capture.Handler{
		Service:   &service,
		Reader:    reader,
		Responder: responder,
	}

	bastion = app.New(ds, []app.Router{&handler}).Bastion
	db = ds.DB()
	code := m.Run()
	clearCollection()
	os.Exit(code)
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
			clearCollection()

			e := gobastion.Tester(t, bastion)
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
			e := gobastion.Tester(t, bastion)
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

	e := gobastion.Tester(t, bastion)
	e.POST("/captures/").
		WithJSON("{").
		Expect().
		Status(http.StatusBadRequest).
		JSON().Object().Equal(response)
}

func TestListCapturesWhenEmpty(t *testing.T) {
	clearCollection()
	e := gobastion.Tester(t, bastion)
	e.GET("/captures/").Expect().JSON().Array().Empty()
}

func TestListCapturesWithValues(t *testing.T) {
	clearCollection()
	if err := createCapture(); err != nil {
		log.Fatal(err)
	}

	e := gobastion.Tester(t, bastion)
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

func createCapture() error {
	c := getCapture(1, 1, "1989-12-26T06:01:00.00Z", []float64{})
	now := time.Now()
	c.CreatedDate, c.LastModified = now, now

	return db.C(capture.Domain).Insert(c)
}
