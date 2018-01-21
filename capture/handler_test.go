package capture_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"os"

	"github.com/ifreddyrondon/gobastion"
	"github.com/ifreddyrondon/gocapture/app"
	"github.com/ifreddyrondon/gocapture/capture"
	"github.com/ifreddyrondon/gocapture/database"
	"gopkg.in/mgo.v2"
)

var application *app.App
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

	handler := new(capture.Handler)
	handler.Reader = reader
	handler.Responder = responder

	application = app.New(ds, []app.Router{handler})
	db = ds.DB()
	code := m.Run()
	clearCollection()
	if err != nil {
		log.Panic(err)
	}

	os.Exit(code)
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	application.Bastion.APIRouter.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func checkErrorResponse(t *testing.T, expected, actual map[string]interface{}) {
	if actual["error"] != expected["error"] {
		t.Errorf("Expected the Error %v. Got '%v'", expected["error"], actual["error"])
	}

	if actual["message"] != expected["message"] {
		t.Errorf("Expected the Message '%v'. Got '%v'", expected["message"], actual["message"])
	}

	if actual["status"] != expected["status"] {
		t.Errorf("Expected the Status '%v'. Got '%v'", expected["status"], actual["status"])
	}
}

// TODO: create test for fields with shadow name and 400
func TestCreateCapture(t *testing.T) {
	tt := []struct {
		name     string
		payload  []byte
		status   int
		response map[string]interface{}
	}{
		{
			name:    "create capture with date",
			payload: []byte(`{"lat": 1, "lng": 12, "date": "1989-12-26T06:01:00.00Z"}`),
			status:  http.StatusCreated,
			response: map[string]interface{}{
				"payload":   "",
				"lat":       1.0,
				"lng":       12.0,
				"timestamp": "1989-12-26T06:01:00Z",
			},
		},
		{
			name:    "create capture with timestamp",
			payload: []byte(`{"lat": 1, "lng": 12, "timestamp": "630655260"}`),
			status:  http.StatusCreated,
			response: map[string]interface{}{
				"payload":   "",
				"lat":       1.0,
				"lng":       12.0,
				"timestamp": "1989-12-26T06:01:00Z",
			},
		},
		{
			name:    "create capture with latitude, longitude and data names",
			payload: []byte(`{"latitude": 1, "longitude": 12, "date": "630655260"}`),
			status:  http.StatusCreated,
			response: map[string]interface{}{
				"payload":   "",
				"lat":       1.0,
				"lng":       12.0,
				"timestamp": "1989-12-26T06:01:00Z",
			},
		},
		{
			name:    "bad request, missing body",
			payload: []byte(`{`),
			status:  http.StatusBadRequest,
			response: map[string]interface{}{
				"status":  400.0,
				"error":   "Bad Request",
				"message": "unexpected EOF",
			},
		},
		{
			name:    "bad request, missing lng",
			payload: []byte(`{"lat": 1, "date": "1989-12-26T06:01:00.00Z"}`),
			status:  http.StatusBadRequest,
			response: map[string]interface{}{
				"status":  400.0,
				"error":   "Bad Request",
				"message": "missing longitude",
			},
		},
		{
			name:    "bad request, missing lat",
			payload: []byte(`{"lng": 12, "date": "1989-12-26T06:01:00.00Z"}`),
			status:  http.StatusBadRequest,
			response: map[string]interface{}{
				"status":  400.0,
				"error":   "Bad Request",
				"message": "missing latitude",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			clearCollection()

			req, _ := http.NewRequest("POST", "/captures/", bytes.NewBuffer(tc.payload))
			response := executeRequest(req)

			checkResponseCode(t, tc.status, response.Code)

			var m map[string]interface{}
			json.Unmarshal(response.Body.Bytes(), &m)

			if tc.status != http.StatusOK {
				checkErrorResponse(t, tc.response, m)
			}

			if m["id"] == "" {
				t.Errorf("Expected id diferent from empty")
			}

			if m["created_date"] == "" {
				t.Errorf("Expected id diferent from empty")
			}

			if m["last_modified"] == "" {
				t.Errorf("Expected id diferent from empty")
			}

			if m["lat"] != tc.response["lat"] {
				t.Errorf("Expected lat to be '%v'. Got '%v'", tc.response["lat"], m["lat"])
			}

			if m["lng"] != tc.response["lng"] {
				t.Errorf("Expected lng to be '%v'. Got '%v'", tc.response["lng"], m["lng"])
			}

			if m["timestamp"] != tc.response["timestamp"] {
				t.Errorf("Expected timestamp to be '%v'. Got '%v'", tc.response["timestamp"], m["timestamp"])
			}
		})
	}
}
