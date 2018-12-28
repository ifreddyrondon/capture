package capture_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/araddon/dateparse"
	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/capture/config"
	"github.com/ifreddyrondon/capture/pkg/capture"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func setup(t *testing.T) (*bastion.Bastion, func()) {
	cfg, err := config.FromString(`PG="postgres://localhost/captures_app_test?sslmode=disable"`)
	if err != nil {
		t.Fatal(err)
	}

	db := cfg.Resources.Get("database").(*gorm.DB)
	store := capture.NewPGStore(db.Table("captures"))
	store.Migrate()

	app := bastion.New()
	app.APIRouter.Mount("/", capture.Routes(store))

	return app, func() { store.Drop() }
}

//func TestBulkCreateValidCapture(t *testing.T) {
//	app, teardown := setup(t)
//	defer teardown()
//
//	e := bastion.Tester(t, app)
//	tt := []struct {
//		name     string
//		payload  []map[string]interface{}
//		response []map[string]interface{}
//	}{
//		{
//			name: "create capture with payload, date and point",
//			payload: []map[string]interface{}{
//				{
//					"latitude":  1,
//					"longitude": 12,
//					"timestamp": "630655260",
//					"payload": []map[string]interface{}{
//						{
//							"name":  "power",
//							"value": []interface{}{-70.0, -100.1, 3.1},
//						},
//					},
//				},
//				{
//					"latitude":  2,
//					"longitude": 3,
//					"timestamp": "630655260",
//					"payload": []map[string]interface{}{
//						{
//							"name":  "power",
//							"value": []interface{}{-45.0, -32.1, 34.1},
//						},
//					},
//				},
//			},
//			response: []map[string]interface{}{
//				{
//					"lat":       1.0,
//					"lng":       12.0,
//					"timestamp": "1989-12-26T06:01:00Z",
//					"tags":      []string{},
//					"payload": []map[string]interface{}{
//						{
//							"name":  "power",
//							"value": []interface{}{-70.0, -100.1, 3.1},
//						},
//					},
//				},
//				{
//					"lat":       2.0,
//					"lng":       3.0,
//					"timestamp": "1989-12-26T06:01:00Z",
//					"tags":      []string{},
//					"payload": []map[string]interface{}{
//						{
//							"name":  "power",
//							"value": []interface{}{-45.0, -32.1, 34.1},
//						},
//					},
//				},
//			},
//		},
//		{
//			name: "create capture with payload and date without point",
//			payload: []map[string]interface{}{
//				{
//					"date": "630655260",
//					"payload": []map[string]interface{}{
//						{
//							"name":  "power",
//							"value": []interface{}{-70.0, -100.1, 3.1},
//						},
//					},
//				},
//				{
//					"date": "630655260",
//					"payload": []map[string]interface{}{
//						{
//							"name":  "power",
//							"value": []interface{}{-50.0, -30.1, 10.1},
//						},
//					},
//				},
//			},
//			response: []map[string]interface{}{
//				{
//					"lat":       nil,
//					"lng":       nil,
//					"timestamp": "1989-12-26T06:01:00Z",
//					"tags":      []string{},
//					"payload": []map[string]interface{}{
//						{
//							"name":  "power",
//							"value": []interface{}{-70.0, -100.1, 3.1},
//						},
//					},
//				},
//				{
//					"lat":       nil,
//					"lng":       nil,
//					"timestamp": "1989-12-26T06:01:00Z",
//					"tags":      []string{},
//					"payload": []map[string]interface{}{
//						{
//							"name":  "power",
//							"value": []interface{}{-50.0, -30.1, 10.1},
//						},
//					},
//				},
//			},
//		},
//		{
//			name: "payload with three capture but one is invalid",
//			payload: []map[string]interface{}{
//				{
//					"date": "630655260",
//					"payload": []map[string]interface{}{
//						{
//							"name":  "power",
//							"value": []interface{}{-70.0, -100.1, 3.1},
//						},
//					},
//				},
//				{
//					"date": "630655260",
//					"payload": []map[string]interface{}{
//						{
//							"name":  "power",
//							"value": []interface{}{-50.0, -30.1, 10.1},
//						},
//					},
//				},
//				{
//					"lat":  -10001.0,
//					"lng":  12.0,
//					"date": "630655260",
//					"payload": []map[string]interface{}{
//						{
//							"name":  "power",
//							"value": []interface{}{-50.0, -30.1, 10.1},
//						},
//					},
//				},
//			},
//			response: []map[string]interface{}{
//				{
//					"lat":       nil,
//					"lng":       nil,
//					"timestamp": "1989-12-26T06:01:00Z",
//					"tags":      []string{},
//					"payload": []map[string]interface{}{
//						{
//							"name":  "power",
//							"value": []interface{}{-70.0, -100.1, 3.1},
//						},
//					},
//				},
//				{
//					"lat":       nil,
//					"lng":       nil,
//					"timestamp": "1989-12-26T06:01:00Z",
//					"tags":      []string{},
//					"payload": []map[string]interface{}{
//						{
//							"name":  "power",
//							"value": []interface{}{-50.0, -30.1, 10.1},
//						},
//					},
//				},
//			},
//		},
//	}
//
//	for _, tc := range tt {
//		t.Run(tc.name, func(t *testing.T) {
//			array := e.POST("/test").
//				WithJSON(tc.payload).
//				Expect().
//				Status(http.StatusCreated).
//				JSON().Array().NotEmpty()
//
//			array.Length().Equal(len(tc.response))
//			for n, val := range array.Iter() {
//				val.Object().
//					ContainsKey("payload").ValueEqual("payload", tc.response[n]["payload"]).
//					ContainsKey("lat").ValueEqual("lat", tc.response[n]["lat"]).
//					ContainsKey("lng").ValueEqual("lng", tc.response[n]["lng"]).
//					ContainsKey("timestamp").ValueEqual("timestamp", tc.response[n]["timestamp"]).
//					ContainsKey("tags").ValueEqual("tags", tc.response[n]["tags"]).
//					ContainsKey("id").NotEmpty().
//					ContainsKey("createdAt").NotEmpty().
//					ContainsKey("updatedAt").NotEmpty()
//			}
//		})
//	}
//}
//
//func TestBulkCreateOnlyOneValidCaptureItShouldReturnObject(t *testing.T) {
//	app, teardown := setup(t)
//	defer teardown()
//
//	e := bastion.Tester(t, app)
//	payload := []map[string]interface{}{
//		{
//			"date": "630655260",
//			"payload": []map[string]interface{}{
//				{
//					"name":  "power",
//					"value": []interface{}{-70.0, -100.1, 3.1},
//				},
//			},
//		},
//		{
//			"lat":  -10001.0,
//			"lng":  12.0,
//			"date": "630655260",
//			"payload": []map[string]interface{}{
//				{
//					"name":  "power",
//					"value": []interface{}{-50.0, -30.1, 10.1},
//				},
//			},
//		},
//	}
//	response := map[string]interface{}{
//		"lat":       nil,
//		"lng":       nil,
//		"timestamp": "1989-12-26T06:01:00Z",
//		"tags":      []string{},
//		"payload": []map[string]interface{}{
//			{
//				"name":  "power",
//				"value": []interface{}{-70.0, -100.1, 3.1},
//			},
//		},
//	}
//
//	e.POST("/test").
//		WithJSON(payload).
//		Expect().
//		Status(http.StatusCreated).
//		JSON().Object().
//		ContainsKey("payload").ValueEqual("payload", response["payload"]).
//		ContainsKey("lat").ValueEqual("lat", response["lat"]).
//		ContainsKey("lng").ValueEqual("lng", response["lng"]).
//		ContainsKey("timestamp").ValueEqual("timestamp", response["timestamp"]).
//		ContainsKey("tags").ValueEqual("tags", response["tags"]).
//		ContainsKey("id").NotEmpty().
//		ContainsKey("createdAt").NotEmpty().
//		ContainsKey("updatedAt").NotEmpty()
//}
//
//func TestBulkCreateInValidCapturesItShouldReturnError(t *testing.T) {
//	app, teardown := setup(t)
//	defer teardown()
//
//	e := bastion.Tester(t, app)
//	tt := []struct {
//		name     string
//		payload  []map[string]interface{}
//		response map[string]interface{}
//	}{
//		{
//			name: "return error if all the captures are invalid",
//			payload: []map[string]interface{}{
//				{
//					"date": "630655260",
//				},
//				{
//					"lat":  -10001.0,
//					"lng":  12.0,
//					"date": "630655260",
//					"payload": []map[string]interface{}{
//						{
//							"name":  "power",
//							"value": []interface{}{-50.0, -30.1, 10.1},
//						},
//					},
//				},
//			},
//			response: map[string]interface{}{
//				"status":  400.0,
//				"error":   "Bad Request",
//				"message": "cannot unmarshal json into valid captures, it needs at least one valid capture",
//			},
//		},
//		{
//			name: "return error if payload contains more than 100 captures",
//			payload: func() []map[string]interface{} {
//				payload := make([]map[string]interface{}, 101)
//				for i := 0; i < 101; i++ {
//					payload = append(payload, randomCapturePayload())
//				}
//				return payload
//			}(),
//			response: map[string]interface{}{
//				"status":  400.0,
//				"error":   "Bad Request",
//				"message": "limited to 100 calls in a single batch request. If it needs to make more calls than that, use multiple batch requests",
//			},
//		},
//	}
//
//	for _, tc := range tt {
//		t.Run(tc.name, func(t *testing.T) {
//			e.POST("/test").
//				WithJSON(tc.payload).
//				Expect().
//				Status(http.StatusBadRequest).
//				JSON().Object().Equal(tc.response)
//		})
//	}
//}

func TestListCapturesWhenEmpty(t *testing.T) {
	app, teardown := setup(t)
	defer teardown()

	e := bastion.Tester(t, app)
	e.GET("/").Expect().JSON().Array().Empty()
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
	e.POST("/test").WithJSON(payload).Expect().Status(http.StatusCreated)

	array := e.GET("/").
		Expect().
		Status(http.StatusOK).
		JSON().Array().NotEmpty()

	array.Length().Equal(1)
	array.First().Object().
		ContainsKey("payload").
		ContainsKey("location").
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
		"location": map[string]float64{
			"lat":       1,
			"lng":       12,
			"elevation": 50,
		},
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
	obj := e.POST("/test").WithJSON(capPayload).Expect().Status(http.StatusCreated).
		JSON().Object().Raw()

	e.GET(fmt.Sprintf("/%v", obj["id"])).Expect().
		Status(http.StatusOK).
		JSON().Object().
		ContainsKey("payload").ValueEqual("payload", capPayload["payload"]).
		ContainsKey("location").ValueEqual("location", capPayload["location"]).
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
	e.GET("/00000000-0000-0000-0000-000000000000").Expect().
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
	e.GET("/ads").Expect().
		Status(http.StatusBadRequest).
		JSON().Object().Equal(response)
}

func TestDeleteCapture(t *testing.T) {
	app, teardown := setup(t)
	defer teardown()

	capPayload := map[string]interface{}{
		"payload": []map[string]interface{}{
			{
				"name":  "power",
				"value": []interface{}{-70.0, -100.1, 3.1},
			},
		},
	}
	e := bastion.Tester(t, app)
	obj := e.POST("/test").WithJSON(capPayload).Expect().Status(http.StatusCreated).
		JSON().Object().Raw()

	id := obj["id"]

	e.GET(fmt.Sprintf("/%v", id)).Expect().Status(http.StatusOK)
	e.DELETE(fmt.Sprintf("/%v", id)).Expect().Status(http.StatusNoContent)
	e.GET(fmt.Sprintf("/%v", id)).Expect().Status(http.StatusNotFound)
}

func TestUpdateCapture(t *testing.T) {
	app, teardown := setup(t)
	defer teardown()

	e := bastion.Tester(t, app)
	capPayload := map[string]interface{}{
		"location": map[string]float64{
			"lat":       1,
			"lng":       12,
			"elevation": 50,
		},
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
				"location": map[string]float64{
					"lat":       1,
					"lng":       capPayload["location"].(map[string]float64)["lng"],
					"elevation": capPayload["location"].(map[string]float64)["elevation"],
				},
				"elevation": capPayload["elevation"],
				"timestamp": capPayload["timestamp"],
				"tags":      capPayload["tags"],
				"payload":   capPayload["payload"],
			},
		},
		{
			"update lng",
			map[string]interface{}{
				"location": map[string]float64{
					"lat":       capPayload["location"].(map[string]float64)["lat"],
					"lng":       30,
					"elevation": capPayload["location"].(map[string]float64)["elevation"],
				},
				"timestamp": capPayload["timestamp"],
				"tags":      capPayload["tags"],
				"payload":   capPayload["payload"],
			},
		},
		{
			"update elevation",
			map[string]interface{}{
				"location": map[string]float64{
					"lat":       capPayload["location"].(map[string]float64)["lat"],
					"lng":       capPayload["location"].(map[string]float64)["lng"],
					"elevation": 100,
				},
				"timestamp": capPayload["timestamp"],
				"tags":      capPayload["tags"],
				"payload":   capPayload["payload"],
			},
		},
		{
			"update timestamp",
			map[string]interface{}{
				"location":  capPayload["location"],
				"timestamp": "2006-07-12T06:01:00Z",
				"tags":      capPayload["tags"],
				"payload":   capPayload["payload"],
			},
		},
		{
			"update payload",
			map[string]interface{}{
				"location":  capPayload["location"],
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
				"location":  capPayload["location"],
				"timestamp": capPayload["timestamp"],
				"payload":   capPayload["payload"],
				"tags":      []string{"tag1", "tag2"},
			},
		},
		{
			"do not update id",
			map[string]interface{}{
				"id":        "123",
				"location":  capPayload["location"],
				"timestamp": capPayload["timestamp"],
				"tags":      capPayload["tags"],
				"payload":   capPayload["payload"],
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			createdObj := e.POST("/test").WithJSON(capPayload).Expect().
				Status(http.StatusCreated).JSON().Object().Raw()

			e.GET(fmt.Sprintf("/%v", createdObj["id"])).Expect().Status(http.StatusOK)
			tc.updatePayload["id"] = createdObj["id"]

			updatedObj := e.PUT(fmt.Sprintf("/%v", createdObj["id"])).WithJSON(tc.updatePayload).Expect().
				Status(http.StatusOK).
				JSON().Object().
				ContainsKey("id").ValueEqual("id", tc.updatePayload["id"]).
				ContainsKey("location").ValueEqual("location", tc.updatePayload["location"]).
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
	capPayload := map[string]interface{}{
		"location": map[string]float64{
			"lat": 1,
			"lng": 12,
		},
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
				"location": map[string]float64{
					"lat": 200,
					"lng": capPayload["location"].(map[string]float64)["lng"],
				},
				"timestamp": capPayload["timestamp"],
				"payload":   capPayload["payload"],
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
				"location": map[string]float64{
					"lat": capPayload["location"].(map[string]float64)["lat"],
					"lng": capPayload["location"].(map[string]float64)["lng"],
				},
				"timestamp": capPayload["timestamp"],
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
			createdObj := e.POST("/test").WithJSON(capPayload).Expect().
				Status(http.StatusCreated).JSON().Object().Raw()
			e.GET(fmt.Sprintf("/%v", createdObj["id"])).Expect().Status(http.StatusOK)
			tc.updatePayload["id"] = createdObj["id"]
			e.PUT(fmt.Sprintf("/%v", createdObj["id"])).WithJSON(tc.updatePayload).Expect().
				Status(http.StatusBadRequest).
				JSON().Object().Equal(tc.response)
		})
	}
}

func TestUpdateCaptureFailsBadRequestWhenMissingID(t *testing.T) {
	app, teardown := setup(t)
	defer teardown()

	e := bastion.Tester(t, app)
	capPayload := map[string]interface{}{
		"location": map[string]float64{
			"lat": 1,
			"lng": 12,
		},
		"timestamp": "1989-12-26T06:01:00Z",
		"payload": []map[string]interface{}{
			{
				"name":  "power",
				"value": []interface{}{-70.0, -100.1, 3.1},
			},
		},
	}

	response := map[string]interface{}{
		"status":  400.0,
		"error":   "Bad Request",
		"message": "capture id must not be blank",
	}

	createdObj := e.POST("/test").WithJSON(capPayload).Expect().
		Status(http.StatusCreated).JSON().Object().Raw()
	e.GET(fmt.Sprintf("/%v", createdObj["id"])).Expect().Status(http.StatusOK)

	e.PUT(fmt.Sprintf("/%v", createdObj["id"])).WithJSON(capPayload).Expect().
		Status(http.StatusBadRequest).
		JSON().Object().Equal(response)
}

//func TestUnmarshalCapturesFail(t *testing.T) {
//	app, teardown := setup(t)
//	defer teardown()
//
//	e := bastion.Tester(t, app)
//	tt := []struct {
//		name     string
//		payload  []interface{}
//		response map[string]interface{}
//	}{
//		{
//			name:    "bad request, missing body",
//			payload: []interface{}{},
//			response: map[string]interface{}{
//				"status":  400.0,
//				"error":   "Bad Request",
//				"message": "cannot unmarshal json into valid captures, it needs at least one valid capture",
//			},
//		},
//	}
//
//	for _, tc := range tt {
//		t.Run(tc.name, func(t *testing.T) {
//			e.POST("/test").
//				WithJSON(tc.payload).
//				Expect().
//				Status(http.StatusBadRequest).
//				JSON().Object().Equal(tc.response)
//		})
//	}
//}
