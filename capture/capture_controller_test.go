package capture_test

import (
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"sync"
	"testing"

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
	repo := capture.NewPGRepository(getDB().Table("captures"))
	repo.Migrate()
	teardown := func() { repo.Drop() }

	controller := capture.NewController(repo, json.NewRender, app.ContextKey("capture"))
	app := bastion.New(bastion.Options{})
	app.APIRouter.Mount("/captures/", controller.Router())

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
				"payload": []map[string]interface{}{
					{
						"name":  "power",
						"value": []interface{}{-70.0, -100.1, 3.1},
					},
				},
			},
			response: map[string]interface{}{
				"lat":       1.0,
				"lng":       12.0,
				"timestamp": "1989-12-26T06:01:00Z",
				"tags":      []string{},
				"payload": []map[string]interface{}{
					{
						"name":  "power",
						"value": []interface{}{-70.0, -100.1, 3.1},
					},
				},
			},
		},
		{
			name: "create capture with multiple metrics in payload",
			payload: map[string]interface{}{
				"latitude":  1,
				"longitude": 12,
				"timestamp": "630655260",
				"payload": []map[string]interface{}{
					{
						"name":  "power",
						"value": []interface{}{-70.0, -100.1, 3.1},
					},
					{
						"name":  "frequency",
						"value": []interface{}{100.0, 200.0, 300.0},
					},
				},
			},
			response: map[string]interface{}{
				"lat":       1.0,
				"lng":       12.0,
				"timestamp": "1989-12-26T06:01:00Z",
				"tags":      []string{},
				"payload": []map[string]interface{}{
					{
						"name":  "power",
						"value": []interface{}{-70.0, -100.1, 3.1},
					},
					{
						"name":  "frequency",
						"value": []interface{}{100.0, 200.0, 300.0},
					},
				},
			},
		},
		{
			name: "create capture with payload, date and point with altitude",
			payload: map[string]interface{}{
				"latitude":  1,
				"longitude": 12,
				"altitude":  50,
				"timestamp": "630655260",
				"payload": []map[string]interface{}{
					{
						"name":  "power",
						"value": []interface{}{-70.0, -100.1, 3.1},
					},
				},
			},
			response: map[string]interface{}{
				"lat":       1.0,
				"lng":       12.0,
				"elevation": 50,
				"timestamp": "1989-12-26T06:01:00Z",
				"tags":      []string{},
				"payload": []map[string]interface{}{
					{
						"name":  "power",
						"value": []interface{}{-70.0, -100.1, 3.1},
					},
				},
			},
		},
		{
			name: "create capture with payload and date without point",
			payload: map[string]interface{}{
				"date": "630655260",
				"payload": []map[string]interface{}{
					{
						"name":  "power",
						"value": []interface{}{-70.0, -100.1, 3.1},
					},
				},
			},
			response: map[string]interface{}{
				"lat":       nil,
				"lng":       nil,
				"timestamp": "1989-12-26T06:01:00Z",
				"tags":      []string{},
				"payload": []map[string]interface{}{
					{
						"name":  "power",
						"value": []interface{}{-70.0, -100.1, 3.1},
					},
				},
			},
		},
		{
			name: "create capture with payload and tags",
			payload: map[string]interface{}{
				"timestamp": "630655260",
				"tags":      []string{"tag1", "tag2"},
				"payload": []map[string]interface{}{
					{
						"name":  "power",
						"value": []interface{}{-70.0, -100.1, 3.1},
					},
				},
			},
			response: map[string]interface{}{
				"lat":       nil,
				"lng":       nil,
				"timestamp": "1989-12-26T06:01:00Z",
				"tags":      []string{"tag1", "tag2"},
				"payload": []map[string]interface{}{
					{
						"name":  "power",
						"value": []interface{}{-70.0, -100.1, 3.1},
					},
				},
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
				ContainsKey("elevation").ValueEqual("elevation", tc.response["elevation"]).
				ContainsKey("timestamp").ValueEqual("timestamp", tc.response["timestamp"]).
				ContainsKey("tags").ValueEqual("tags", tc.response["tags"]).
				ContainsKey("id").NotEmpty().
				ContainsKey("createdAt").NotEmpty().
				ContainsKey("updatedAt").NotEmpty()
		})
	}
}

func TestCreateOnlyPayloadCapture(t *testing.T) {
	app, teardown := setup(t)
	defer teardown()

	body := map[string]interface{}{
		"payload": []map[string]interface{}{
			{
				"name":  "power",
				"value": []interface{}{-70.0, -100.1, 3.1},
			},
		},
	}
	e := bastion.Tester(t, app)
	e.POST("/captures/").
		WithJSON(body).
		Expect().
		Status(http.StatusCreated).
		JSON().Object().
		ContainsKey("payload").ValueEqual("payload", body["payload"]).
		ContainsKey("tags").ValueEqual("tags", []string{}).
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
				"message": "payload value must not be blank",
			},
		},
		{
			name:    "bad request, missing lng",
			payload: map[string]interface{}{"lat": 1, "date": "630655260"},
			response: map[string]interface{}{
				"status":  400.0,
				"error":   "Bad Request",
				"message": "longitude must not be blank",
			},
		},
		{
			name:    "bad request, missing lat",
			payload: map[string]interface{}{"lng": 1, "date": "630655260"},
			response: map[string]interface{}{
				"status":  400.0,
				"error":   "Bad Request",
				"message": "latitude must not be blank",
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

func TestBulkCreateValidCapture(t *testing.T) {
	app, teardown := setup(t)
	defer teardown()

	e := bastion.Tester(t, app)
	tt := []struct {
		name     string
		payload  []map[string]interface{}
		response []map[string]interface{}
	}{
		{
			name: "create capture with payload, date and point",
			payload: []map[string]interface{}{
				{
					"latitude":  1,
					"longitude": 12,
					"timestamp": "630655260",
					"payload": []map[string]interface{}{
						{
							"name":  "power",
							"value": []interface{}{-70.0, -100.1, 3.1},
						},
					},
				},
				{
					"latitude":  2,
					"longitude": 3,
					"timestamp": "630655260",
					"payload": []map[string]interface{}{
						{
							"name":  "power",
							"value": []interface{}{-45.0, -32.1, 34.1},
						},
					},
				},
			},
			response: []map[string]interface{}{
				{
					"lat":       1.0,
					"lng":       12.0,
					"timestamp": "1989-12-26T06:01:00Z",
					"tags":      []string{},
					"payload": []map[string]interface{}{
						{
							"name":  "power",
							"value": []interface{}{-70.0, -100.1, 3.1},
						},
					},
				},
				{
					"lat":       2.0,
					"lng":       3.0,
					"timestamp": "1989-12-26T06:01:00Z",
					"tags":      []string{},
					"payload": []map[string]interface{}{
						{
							"name":  "power",
							"value": []interface{}{-45.0, -32.1, 34.1},
						},
					},
				},
			},
		},
		{
			name: "create capture with payload and date without point",
			payload: []map[string]interface{}{
				{
					"date": "630655260",
					"payload": []map[string]interface{}{
						{
							"name":  "power",
							"value": []interface{}{-70.0, -100.1, 3.1},
						},
					},
				},
				{
					"date": "630655260",
					"payload": []map[string]interface{}{
						{
							"name":  "power",
							"value": []interface{}{-50.0, -30.1, 10.1},
						},
					},
				},
			},
			response: []map[string]interface{}{
				{
					"lat":       nil,
					"lng":       nil,
					"timestamp": "1989-12-26T06:01:00Z",
					"tags":      []string{},
					"payload": []map[string]interface{}{
						{
							"name":  "power",
							"value": []interface{}{-70.0, -100.1, 3.1},
						},
					},
				},
				{
					"lat":       nil,
					"lng":       nil,
					"timestamp": "1989-12-26T06:01:00Z",
					"tags":      []string{},
					"payload": []map[string]interface{}{
						{
							"name":  "power",
							"value": []interface{}{-50.0, -30.1, 10.1},
						},
					},
				},
			},
		},
		{
			name: "payload with three capture but one is invalid",
			payload: []map[string]interface{}{
				{
					"date": "630655260",
					"payload": []map[string]interface{}{
						{
							"name":  "power",
							"value": []interface{}{-70.0, -100.1, 3.1},
						},
					},
				},
				{
					"date": "630655260",
					"payload": []map[string]interface{}{
						{
							"name":  "power",
							"value": []interface{}{-50.0, -30.1, 10.1},
						},
					},
				},
				{
					"lat":  -10001.0,
					"lng":  12.0,
					"date": "630655260",
					"payload": []map[string]interface{}{
						{
							"name":  "power",
							"value": []interface{}{-50.0, -30.1, 10.1},
						},
					},
				},
			},
			response: []map[string]interface{}{
				{
					"lat":       nil,
					"lng":       nil,
					"timestamp": "1989-12-26T06:01:00Z",
					"tags":      []string{},
					"payload": []map[string]interface{}{
						{
							"name":  "power",
							"value": []interface{}{-70.0, -100.1, 3.1},
						},
					},
				},
				{
					"lat":       nil,
					"lng":       nil,
					"timestamp": "1989-12-26T06:01:00Z",
					"tags":      []string{},
					"payload": []map[string]interface{}{
						{
							"name":  "power",
							"value": []interface{}{-50.0, -30.1, 10.1},
						},
					},
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			array := e.POST("/captures/").
				WithJSON(tc.payload).
				Expect().
				Status(http.StatusCreated).
				JSON().Array().NotEmpty()

			array.Length().Equal(len(tc.response))
			for n, val := range array.Iter() {
				val.Object().
					ContainsKey("payload").ValueEqual("payload", tc.response[n]["payload"]).
					ContainsKey("lat").ValueEqual("lat", tc.response[n]["lat"]).
					ContainsKey("lng").ValueEqual("lng", tc.response[n]["lng"]).
					ContainsKey("timestamp").ValueEqual("timestamp", tc.response[n]["timestamp"]).
					ContainsKey("tags").ValueEqual("tags", tc.response[n]["tags"]).
					ContainsKey("id").NotEmpty().
					ContainsKey("createdAt").NotEmpty().
					ContainsKey("updatedAt").NotEmpty()
			}
		})
	}
}

func TestBulkCreateOnlyOneValidCaptureItShouldReturnObject(t *testing.T) {
	app, teardown := setup(t)
	defer teardown()

	e := bastion.Tester(t, app)
	payload := []map[string]interface{}{
		{
			"date": "630655260",
			"payload": []map[string]interface{}{
				{
					"name":  "power",
					"value": []interface{}{-70.0, -100.1, 3.1},
				},
			},
		},
		{
			"lat":  -10001.0,
			"lng":  12.0,
			"date": "630655260",
			"payload": []map[string]interface{}{
				{
					"name":  "power",
					"value": []interface{}{-50.0, -30.1, 10.1},
				},
			},
		},
	}
	response := map[string]interface{}{
		"lat":       nil,
		"lng":       nil,
		"timestamp": "1989-12-26T06:01:00Z",
		"tags":      []string{},
		"payload": []map[string]interface{}{
			{
				"name":  "power",
				"value": []interface{}{-70.0, -100.1, 3.1},
			},
		},
	}

	e.POST("/captures/").
		WithJSON(payload).
		Expect().
		Status(http.StatusCreated).
		JSON().Object().
		ContainsKey("payload").ValueEqual("payload", response["payload"]).
		ContainsKey("lat").ValueEqual("lat", response["lat"]).
		ContainsKey("lng").ValueEqual("lng", response["lng"]).
		ContainsKey("timestamp").ValueEqual("timestamp", response["timestamp"]).
		ContainsKey("tags").ValueEqual("tags", response["tags"]).
		ContainsKey("id").NotEmpty().
		ContainsKey("createdAt").NotEmpty().
		ContainsKey("updatedAt").NotEmpty()
}

func TestBulkCreateInValidCapturesItShouldReturnError(t *testing.T) {
	app, teardown := setup(t)
	defer teardown()

	e := bastion.Tester(t, app)
	tt := []struct {
		name     string
		payload  []map[string]interface{}
		response map[string]interface{}
	}{
		{
			name: "return error if all the captures are invalid",
			payload: []map[string]interface{}{
				{
					"date": "630655260",
				},
				{
					"lat":  -10001.0,
					"lng":  12.0,
					"date": "630655260",
					"payload": []map[string]interface{}{
						{
							"name":  "power",
							"value": []interface{}{-50.0, -30.1, 10.1},
						},
					},
				},
			},
			response: map[string]interface{}{
				"status":  400.0,
				"error":   "Bad Request",
				"message": "cannot unmarshal json into valid captures, it needs at least one valid capture",
			},
		},
		{
			name: "return error if payload contains more than 100 captures",
			payload: func() []map[string]interface{} {
				payload := make([]map[string]interface{}, 101)
				for i := 0; i < 101; i++ {
					payload = append(payload, randomCapturePayload())
				}
				return payload
			}(),
			response: map[string]interface{}{
				"status":  400.0,
				"error":   "Bad Request",
				"message": "limited to 100 calls in a single batch request. If it needs to make more calls than that, use multiple batch requests",
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
		"payload": []map[string]interface{}{
			{
				"name":  "power",
				"value": []interface{}{-70.0, -100.1, 3.1},
			},
		},
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
		ContainsKey("elevation").
		ContainsKey("timestamp").
		ContainsKey("id").
		ContainsKey("tags").
		ContainsKey("createdAt").
		ContainsKey("updatedAt")
}

func TestGetCapture(t *testing.T) {
	app, teardown := setup(t)
	defer teardown()

	capPayload := map[string]interface{}{
		"lat":       1.0,
		"lng":       12.0,
		"elevation": 50,
		"timestamp": "1989-12-26T06:01:00Z",
		"tags":      []string{"tag1", "tag2"},
		"payload": []map[string]interface{}{
			{
				"name":  "power",
				"value": []interface{}{-70.0, -100.1, 3.1},
			},
		},
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
		ContainsKey("elevation").ValueEqual("elevation", capPayload["elevation"]).
		ContainsKey("timestamp").NotEmpty().
		ContainsKey("tags").ValueEqual("tags", capPayload["tags"]).
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
	e.GET("/captures/00000000-0000-0000-0000-000000000000").Expect().
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
	e.GET("/captures/ads").Expect().
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
		"payload": []map[string]interface{}{
			{
				"name":  "power",
				"value": []interface{}{-70.0, -100.1, 3.1},
			},
		},
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
		"elevation": 30.0,
		"timestamp": "1989-12-26T06:01:00Z",
		"tags":      []string{},
		"payload": []map[string]interface{}{
			{
				"name":  "power",
				"value": []interface{}{-70.0, -100.1, 3.1},
			},
		},
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
				"elevation": capPayload["elevation"],
				"timestamp": capPayload["timestamp"],
				"tags":      capPayload["tags"],
				"payload":   capPayload["payload"],
			},
		},
		{
			"update lng",
			map[string]interface{}{
				"lat":       capPayload["lat"],
				"lng":       30,
				"elevation": capPayload["elevation"],
				"timestamp": capPayload["timestamp"],
				"tags":      capPayload["tags"],
				"payload":   capPayload["payload"],
			},
		},
		{
			"update elevation",
			map[string]interface{}{
				"lat":       capPayload["lat"],
				"lng":       capPayload["lng"],
				"elevation": 100,
				"timestamp": capPayload["timestamp"],
				"tags":      capPayload["tags"],
				"payload":   capPayload["payload"],
			},
		},
		{
			"update timestamp",
			map[string]interface{}{
				"lat":       capPayload["lat"],
				"lng":       capPayload["lng"],
				"elevation": capPayload["elevation"],
				"timestamp": "2006-07-12T06:01:00Z",
				"tags":      capPayload["tags"],
				"payload":   capPayload["payload"],
			},
		},
		{
			"update payload",
			map[string]interface{}{
				"lat":       capPayload["lat"],
				"lng":       capPayload["lng"],
				"elevation": capPayload["elevation"],
				"timestamp": capPayload["timestamp"],
				"tags":      capPayload["tags"],
				"payload": []map[string]interface{}{
					{
						"name":  "power",
						"value": []interface{}{1},
					},
				},
			},
		},
		{
			"update tags",
			map[string]interface{}{
				"lat":       capPayload["lat"],
				"lng":       capPayload["lng"],
				"elevation": capPayload["elevation"],
				"timestamp": capPayload["timestamp"],
				"payload":   capPayload["payload"],
				"tags":      []string{"tag1", "tag2"},
			},
		},
		{
			"do not update id",
			map[string]interface{}{
				"id":        "123",
				"lat":       capPayload["lat"],
				"lng":       capPayload["lng"],
				"elevation": capPayload["elevation"],
				"timestamp": capPayload["timestamp"],
				"tags":      capPayload["tags"],
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
				ContainsKey("elevation").ValueEqual("elevation", tc.updatePayload["elevation"]).
				ContainsKey("timestamp").ValueEqual("timestamp", tc.updatePayload["timestamp"]).
				ContainsKey("payload").ValueEqual("payload", tc.updatePayload["payload"]).
				ContainsKey("tags").ValueEqual("tags", tc.updatePayload["tags"]).
				ContainsKey("createdAt").NotEmpty().
				ContainsKey("updatedAt").NotEmpty().
				Raw()

			// updatedAt from put should be after updatedAt from post
			updatedAtFromCreate, err := dateparse.ParseAny(createdObj["updatedAt"].(string))
			assert.Nil(t, err)
			updatedAtFromUpdate, err := dateparse.ParseAny(updatedObj["updatedAt"].(string))
			assert.Nil(t, err)
			assert.True(t, updatedAtFromUpdate.After(updatedAtFromCreate))
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
		"payload": []map[string]interface{}{
			{
				"name":  "power",
				"value": []interface{}{-70.0, -100.1, 3.1},
			},
		},
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
				"message": "payload value must not be blank",
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
				"message": "cannot unmarshal json into valid captures, it needs at least one valid capture",
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

func randomCapturePayload() map[string]interface{} {
	return map[string]interface{}{
		"payload": map[string]interface{}{
			"power": []interface{}{
				getRandomPower(),
				getRandomPower(),
				getRandomPower(),
			},
		},
	}
}

func getRandomPower() float64 {
	p := -150 + rand.Float64()*(-10+(-150))
	return math.Ceil(p*100) / 100
}
