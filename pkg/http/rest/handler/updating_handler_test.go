package handler_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/ifreddyrondon/bastion"
	"gopkg.in/src-d/go-kallax.v1"

	"github.com/ifreddyrondon/capture/pkg/domain"
	"github.com/ifreddyrondon/capture/pkg/http/rest/handler"
	"github.com/ifreddyrondon/capture/pkg/updating"
)

type mockUpdatingCaptureService struct {
	err error
}

func (m *mockUpdatingCaptureService) Update(updating.Capture, *domain.Capture) error {
	return m.err
}

func setupUpdatingCaptureHandler(s updating.CaptureService, m func(http.Handler) http.Handler) *bastion.Bastion {
	app := bastion.New()
	app.Use(m)
	app.Put("/", handler.UpdatingCapture(s))
	return app
}

func TestUpdatingCaptureSuccess(t *testing.T) {
	t.Parallel()

	captureID := kallax.NewULID()
	capt := &domain.Capture{
		ID: captureID,
		Payload: domain.Payload{
			domain.Metric{Name: "power", Value: []float64{-70.0, -100.1, 3.1}},
		},
		Timestamp: s2t("1989-12-26T06:01:00.00Z"),
		Location:  &domain.Point{LAT: f2P(1), LNG: f2P(1), Elevation: f2P(1)},
		Tags:      []string{"at night"},
	}

	body := map[string]interface{}{
		"location": map[string]float64{
			"lat":       10,
			"lng":       *capt.Location.LNG,
			"elevation": *capt.Location.Elevation,
		},
	}

	expected := map[string]interface{}{
		"id":        captureID.String(),
		"payload":   []map[string]interface{}{{"name": "power", "value": []interface{}{-70.0, -100.1, 3.1}}},
		"timestamp": "1989-12-26T06:01:00Z",
		"location":  map[string]float64{"lat": 1, "lng": 1, "elevation": 1},
		"tags":      []string{"at night"},
	}

	s := &mockUpdatingCaptureService{}
	app := setupUpdatingCaptureHandler(s, withCaptureMiddle(capt))
	e := bastion.Tester(t, app)

	e.PUT("/").WithJSON(body).Expect().
		Status(http.StatusOK).
		JSON().Object().
		ContainsKey("id").ValueEqual("id", expected["id"]).
		ContainsKey("location").ValueEqual("location", expected["location"]).
		ContainsKey("timestamp").ValueEqual("timestamp", expected["timestamp"]).
		ContainsKey("payload").ValueEqual("payload", expected["payload"]).
		ContainsKey("tags").ValueEqual("tags", expected["tags"]).
		ContainsKey("createdAt").NotEmpty().
		ContainsKey("updatedAt").NotEmpty().
		Raw()
}

func TestUpdatingCaptureFailsGettingCapture(t *testing.T) {
	t.Parallel()

	app := setupUpdatingCaptureHandler(&mockUpdatingCaptureService{}, withCaptureMiddle(nil))

	body := map[string]interface{}{
		"location": map[string]float64{"lat": 10, "lng": 1, "elevation": 1},
	}
	response := map[string]interface{}{
		"status":  500.0,
		"error":   "Internal Server Error",
		"message": "looks like something went wrong",
	}

	e := bastion.Tester(t, app)
	e.PUT("/").WithJSON(body).Expect().
		Status(http.StatusInternalServerError).
		JSON().Object().Equal(response)
}

func TestUpdatingCaptureFailsUpdating(t *testing.T) {
	t.Parallel()
	s := &mockUpdatingCaptureService{err: errors.New("test")}
	app := setupUpdatingCaptureHandler(s, withCaptureMiddle(defaultCapture))

	body := map[string]interface{}{
		"location": map[string]float64{"lat": 10, "lng": 1, "elevation": 1},
	}
	response := map[string]interface{}{
		"status":  500.0,
		"error":   "Internal Server Error",
		"message": "looks like something went wrong",
	}

	e := bastion.Tester(t, app)
	e.PUT("/").WithJSON(body).Expect().
		Status(http.StatusInternalServerError).
		JSON().Object().Equal(response)
}

func TestUpdatingCaptureFailBadRequest(t *testing.T) {
	t.Parallel()
	tt := []struct {
		name     string
		payload  map[string]interface{}
		response map[string]interface{}
	}{
		{
			name: "bad request, missing lng",
			payload: map[string]interface{}{
				"location": map[string]float64{"lat": 1},
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
				"location": map[string]float64{"lat": 1000, "lng": 1},
			},
			response: map[string]interface{}{
				"status":  400.0,
				"error":   "Bad Request",
				"message": "latitude out of boundaries, may range from -90.0 to 90.0",
			},
		},
	}

	s := &mockUpdatingCaptureService{}
	app := setupUpdatingCaptureHandler(s, withCaptureMiddle(defaultCapture))
	e := bastion.Tester(t, app)
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			e.PUT("/").
				WithJSON(tc.payload).
				Expect().
				Status(http.StatusBadRequest).
				JSON().Object().Equal(tc.response)
		})
	}
}
