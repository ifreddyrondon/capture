package updating_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"gopkg.in/src-d/go-kallax.v1"

	"github.com/ifreddyrondon/capture/pkg/domain"
	"github.com/ifreddyrondon/capture/pkg/updating"
	"github.com/ifreddyrondon/capture/pkg/validator"
)

type mockStore struct {
	err error
}

func (m *mockStore) Save(*domain.Capture) error { return m.err }

var (
	defaultCaptureID = kallax.NewULID()
	defaultCapture   = domain.Capture{
		ID: defaultCaptureID,
		Payload: domain.Payload{
			domain.Metric{Name: "power", Value: 10.0},
		},
		Location: &domain.Point{LAT: f2P(1), LNG: f2P(1)},
		Tags:     []string{},
	}
)

func TestServiceUpdateCaptureOK(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name     string
		payl     updating.Capture
		expected domain.Capture
	}{
		{
			name: "given an empty body should return the same capture",
			payl: updating.Capture{
				Payload: &validator.Payload{
					Payload: []domain.Metric{
						{Name: "power", Value: 10.0},
					},
				},
			},
			expected: func() domain.Capture {
				capt := defaultCapture
				return capt
			}(),
		},
		{
			name: "given payload should update only the payload",
			payl: updating.Capture{
				Payload: &validator.Payload{
					Payload: []domain.Metric{
						{Name: "power", Value: []float64{10.0, 20.0}},
						{Name: "frequency", Value: 300.0},
					},
				},
			},
			expected: func() domain.Capture {
				capt := defaultCapture
				capt.Payload = []domain.Metric{
					{Name: "power", Value: []float64{10.0, 20.0}},
					{Name: "frequency", Value: 300.0},
				}
				return capt
			}(),
		},
		{
			name: "given payload and location should update both",
			payl: updating.Capture{
				Payload: &validator.Payload{
					Payload: []domain.Metric{
						{Name: "power", Value: []float64{20.0}},
					},
				},
				Location: &validator.GeoLocation{LAT: f2P(3), LNG: f2P(5)},
			},
			expected: func() domain.Capture {
				capt := defaultCapture
				capt.Payload = []domain.Metric{{Name: "power", Value: []float64{20.0}}}
				capt.Location = &domain.Point{LAT: f2P(3), LNG: f2P(5)}
				return capt
			}(),
		},
		{
			name: "given tags should update them",
			payl: updating.Capture{
				Tags: []string{"at night"},
			},
			expected: func() domain.Capture {
				capt := defaultCapture
				capt.Tags = []string{"at night"}
				return capt
			}(),
		},
		{
			name: "given timestamp should update timestamp",
			payl: updating.Capture{
				Timestamp: &validator.Timestamp{
					Date: s2n("1989-12-26T06:01:00.00Z"),
					Time: s2t("1989-12-26T06:01:00.00Z"),
				},
			},
			expected: func() domain.Capture {
				capt := defaultCapture
				capt.Timestamp = *s2t("1989-12-26T06:01:00.00Z")
				return capt
			}(),
		},
	}

	s := updating.NewCaptureService(&mockStore{})
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			crrTime := time.Now()
			capt := defaultCapture
			err := s.Update(tc.payl, &capt)
			assert.Nil(t, err)

			assert.Equal(t, tc.expected.ID, capt.ID)
			assert.Equal(t, tc.expected.Payload, capt.Payload)
			assert.Equal(t, tc.expected.Location, capt.Location)
			assert.Equal(t, tc.expected.Tags, capt.Tags)
			assert.NotNil(t, capt.CreatedAt)
			assert.True(t, capt.UpdatedAt.After(crrTime))
			assert.Nil(t, capt.DeletedAt)
		})
	}
}

func TestServiceUpdateCaptureErrWhenSaving(t *testing.T) {
	t.Parallel()
	s := updating.NewCaptureService(&mockStore{err: errors.New("test")})
	data := updating.Capture{
		Payload: &validator.Payload{
			Payload: []domain.Metric{
				{Name: "power", Value: 10.0},
			},
		},
	}
	capt := defaultCapture
	err := s.Update(data, &capt)
	assert.EqualError(t, err, fmt.Sprintf("could not update capture %v: test", defaultCaptureID))
}
