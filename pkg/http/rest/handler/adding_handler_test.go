package handler_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/ifreddyrondon/capture/pkg/http/rest/handler"

	"github.com/ifreddyrondon/bastion"

	"github.com/ifreddyrondon/capture/pkg/adding"
	"github.com/ifreddyrondon/capture/pkg/domain"
)

type mockAddingCaptureService struct {
	capt *domain.Capture
	err  error
}

func (m *mockAddingCaptureService) AddCapture(r *domain.Repository, c adding.Capture) (*domain.Capture, error) {
	return m.capt, m.err
}

func setupAddingCaptureHandler(s adding.CaptureService, m func(http.Handler) http.Handler) *bastion.Bastion {
	app := bastion.New()
	app.Use(m)
	app.Post("/", handler.AddingCapture(s))
	return app
}

func TestAddingCaptureSuccess(t *testing.T) {
	t.Parallel()

	capt := &domain.Capture{
		Payload: domain.Payload{
			domain.Metric{Name: "power", Value: []float64{-70.0, -100.1, 3.1}},
		},
		Timestamp: s2t("1989-12-26T06:01:00.00Z"),
		Location:  &domain.Point{LAT: f2P(1), LNG: f2P(1), Elevation: f2P(1)},
		Tags:      []string{"at night"},
	}

	s := &mockAddingCaptureService{capt: capt}
	app := setupAddingCaptureHandler(s, withRepoMiddle(defaultRepo))
	e := bastion.Tester(t, app)

	payload := map[string]interface{}{
		"location": map[string]float64{
			"latitude":  1,
			"longitude": 1,
			"elevation": 1,
		},
		"timestamp": "1989-12-26T06:01:00.00Z",
		"payload": []map[string]interface{}{
			{
				"name":  "power",
				"value": []interface{}{-70.0, -100.1, 3.1},
			},
		},
	}
	response := map[string]interface{}{
		"location": map[string]float64{
			"lat":       1,
			"lng":       1,
			"elevation": 1,
		},
		"timestamp": "1989-12-26T06:01:00Z",
		"tags":      []string{"at night"},
		"payload": []map[string]interface{}{
			{
				"name":  "power",
				"value": []interface{}{-70.0, -100.1, 3.1},
			},
		},
	}

	e.POST("/").
		WithJSON(payload).
		Expect().
		Status(http.StatusCreated).
		JSON().Object().
		ContainsKey("payload").ValueEqual("payload", response["payload"]).
		ContainsKey("location").ValueEqual("location", response["location"]).
		ContainsKey("timestamp").ValueEqual("timestamp", response["timestamp"]).
		ContainsKey("tags").ValueEqual("tags", response["tags"]).
		ContainsKey("id").NotEmpty().
		ContainsKey("createdAt").NotEmpty().
		ContainsKey("updatedAt").NotEmpty()
}

func TestAddingCaptureFailBadRequest(t *testing.T) {
	t.Parallel()
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
			name: "bad request, missing lng",
			payload: map[string]interface{}{
				"payload": []map[string]interface{}{
					{
						"name":  "power",
						"value": []interface{}{-70.0, -100.1, 3.1},
					},
				},
				"location": map[string]float64{
					"lat": 1,
				},
				"date": "630655260",
			},
			response: map[string]interface{}{
				"status":  400.0,
				"error":   "Bad Request",
				"message": "longitude must not be blank",
			},
		},
		{
			name: "bad request, missing lat",
			payload: map[string]interface{}{
				"payload": []map[string]interface{}{
					{
						"name":  "power",
						"value": []interface{}{-70.0, -100.1, 3.1},
					},
				},
				"location": map[string]float64{
					"lng": 1,
				},
				"date": "630655260",
			},
			response: map[string]interface{}{
				"status":  400.0,
				"error":   "Bad Request",
				"message": "latitude must not be blank",
			},
		},
	}

	s := &mockAddingCaptureService{}
	app := setupAddingCaptureHandler(s, withRepoMiddle(defaultRepo))
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

func TestGettingRepoFromAddingCaptureInternalServer(t *testing.T) {
	t.Parallel()

	s := &mockAddingCaptureService{}
	app := setupAddingCaptureHandler(s, withRepoMiddle(nil))

	payload := map[string]interface{}{
		"payload": []map[string]interface{}{
			{
				"name":  "power",
				"value": []interface{}{-70.0, -100.1, 3.1},
			},
		},
	}
	response := map[string]interface{}{
		"status":  500.0,
		"error":   "Internal Server Error",
		"message": "looks like something went wrong",
	}

	e := bastion.Tester(t, app)
	e.POST("/").
		WithJSON(payload).
		Expect().
		Status(http.StatusInternalServerError).
		JSON().Object().Equal(response)
}

func TestAddingCaptureInternalServer(t *testing.T) {
	t.Parallel()

	s := &mockAddingCaptureService{err: errors.New("test")}
	app := setupAddingCaptureHandler(s, withRepoMiddle(defaultRepo))

	payload := map[string]interface{}{
		"payload": []map[string]interface{}{
			{
				"name":  "power",
				"value": []interface{}{-70.0, -100.1, 3.1},
			},
		},
	}
	response := map[string]interface{}{
		"status":  500.0,
		"error":   "Internal Server Error",
		"message": "looks like something went wrong",
	}

	e := bastion.Tester(t, app)
	e.POST("/").
		WithJSON(payload).
		Expect().
		Status(http.StatusInternalServerError).
		JSON().Object().Equal(response)
}

type mockAddingMultiCaptureService struct {
	captures []domain.Capture
	err      error
}

func (m *mockAddingMultiCaptureService) AddCaptures(*domain.Repository, adding.MultiCapture) ([]domain.Capture, error) {
	return m.captures, m.err
}

func setupAddingMultiCaptureHandler(s adding.MultiCaptureService, m func(http.Handler) http.Handler) *bastion.Bastion {
	app := bastion.New()
	app.Use(m)
	app.Post("/", handler.AddingMultiCapture(s))
	return app
}

func TestAddingMultiCaptureSuccess(t *testing.T) {
	t.Parallel()

	captures := []domain.Capture{
		{
			Payload: domain.Payload{domain.Metric{Name: "power", Value: []int{10.0}}},
			Tags:    []string{},
		},
		{
			Payload: domain.Payload{domain.Metric{Name: "power", Value: []int{30.0}}},
			Tags:    []string{},
		},
	}

	s := &mockAddingMultiCaptureService{captures: captures}
	app := setupAddingMultiCaptureHandler(s, withRepoMiddle(defaultRepo))
	e := bastion.Tester(t, app)

	payload := map[string]interface{}{
		"captures": []map[string]interface{}{
			{
				"payload": []map[string]interface{}{
					{"name": "power", "value": []interface{}{10.0}},
				},
			},
			{
				"payload": []map[string]interface{}{
					{"name": "power", "value": []interface{}{30.0}},
				},
			},
		},
	}
	response := []map[string]interface{}{
		{
			"location": nil,
			"tags":     []string{},
			"payload": []map[string]interface{}{
				{"name": "power", "value": []interface{}{10.0}},
			},
		},
		{
			"location": nil,
			"tags":     []string{},
			"payload": []map[string]interface{}{
				{"name": "power", "value": []interface{}{30.0}},
			},
		},
	}

	res := e.POST("/").
		WithJSON(payload).
		Expect().
		Status(http.StatusCreated).
		JSON().Array()

	res.Length().Equal(2)
	for i, v := range res.Iter() {
		v.Object().
			ContainsKey("payload").ValueEqual("payload", response[i]["payload"]).
			ContainsKey("location").ValueEqual("location", response[i]["location"]).
			ContainsKey("tags").ValueEqual("tags", response[i]["tags"]).
			ContainsKey("id").NotEmpty().
			ContainsKey("timestamp").NotEmpty().
			ContainsKey("createdAt").NotEmpty().
			ContainsKey("updatedAt").NotEmpty()
	}
}

func TestAddingMultiCaptureFailBadRequest(t *testing.T) {
	t.Parallel()

	s := &mockAddingMultiCaptureService{}
	app := setupAddingMultiCaptureHandler(s, withRepoMiddle(nil))

	payload := map[string]interface{}{
		"captures": []map[string]interface{}{},
	}
	response := map[string]interface{}{
		"status":  400.0,
		"error":   "Bad Request",
		"message": "captures value must not be blank or empty",
	}

	e := bastion.Tester(t, app)
	e.POST("/").
		WithJSON(payload).
		Expect().
		Status(http.StatusBadRequest).
		JSON().Object().Equal(response)
}

func TestGettingRepoFromAddingMultiCaptureInternalServer(t *testing.T) {
	t.Parallel()

	s := &mockAddingMultiCaptureService{}
	app := setupAddingMultiCaptureHandler(s, withRepoMiddle(nil))

	payload := map[string]interface{}{
		"captures": []map[string]interface{}{
			{
				"payload": []map[string]interface{}{
					{"name": "power", "value": []interface{}{10.0}},
				},
			},
			{
				"payload": []map[string]interface{}{
					{"name": "power", "value": []interface{}{30.0}},
				},
			},
		},
	}
	response := map[string]interface{}{
		"status":  500.0,
		"error":   "Internal Server Error",
		"message": "looks like something went wrong",
	}

	e := bastion.Tester(t, app)
	e.POST("/").
		WithJSON(payload).
		Expect().
		Status(http.StatusInternalServerError).
		JSON().Object().Equal(response)
}

func TestAddingMultiCaptureInternalServer(t *testing.T) {
	t.Parallel()

	s := &mockAddingMultiCaptureService{err: errors.New("test")}
	app := setupAddingMultiCaptureHandler(s, withRepoMiddle(defaultRepo))

	payload := map[string]interface{}{
		"captures": []map[string]interface{}{
			{
				"payload": []map[string]interface{}{
					{"name": "power", "value": []interface{}{10.0}},
				},
			},
			{
				"payload": []map[string]interface{}{
					{"name": "power", "value": []interface{}{30.0}},
				},
			},
		},
	}
	response := map[string]interface{}{
		"status":  500.0,
		"error":   "Internal Server Error",
		"message": "looks like something went wrong",
	}

	e := bastion.Tester(t, app)
	e.POST("/").
		WithJSON(payload).
		Expect().
		Status(http.StatusInternalServerError).
		JSON().Object().Equal(response)
}
