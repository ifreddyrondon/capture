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

type mockAddingService struct {
	capt *domain.Capture
	err  error
}

func (m *mockAddingService) AddCapture(r *domain.Repository, c adding.Capture) (*domain.Capture, error) {
	return m.capt, m.err
}

func setupAddingHandler(s adding.CaptureService, m func(http.Handler) http.Handler) *bastion.Bastion {
	app := bastion.New()
	app.APIRouter.Use(m)
	app.APIRouter.Post("/", handler.AddingCapture(s))
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

	s := &mockAddingService{capt: capt}
	app := setupAddingHandler(s, withRepoMiddle(defaultRepo))
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

	s := &mockAddingService{}
	app := setupAddingHandler(s, withRepoMiddle(defaultRepo))
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

	s := &mockAddingService{}
	app := setupAddingHandler(s, withRepoMiddle(nil))

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

	s := &mockAddingService{err: errors.New("test")}
	app := setupAddingHandler(s, withRepoMiddle(defaultRepo))

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
