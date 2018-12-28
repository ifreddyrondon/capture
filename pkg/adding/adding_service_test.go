package adding_test

import (
	"testing"
	"time"

	"github.com/ifreddyrondon/capture/pkg/adding"
	"github.com/ifreddyrondon/capture/pkg/domain"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"gopkg.in/src-d/go-kallax.v1"
)

type mockStore struct {
	err error
}

func (m *mockStore) CreateCapture(*domain.Capture) error { return m.err }

func TestServiceAddCaptureOKWithDefaultTimestamp(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name     string
		payl     adding.Capture
		expected domain.Capture
	}{
		{
			name: "given a only name payload should return a domain capture without location and tags",
			payl: adding.Capture{
				Payload: adding.Payload{
					Payload: []domain.Metric{
						{Name: "power", Value: 10.0},
					},
				},
			},
			expected: domain.Capture{
				Payload: domain.Payload{
					domain.Metric{Name: "power", Value: 10.0},
				},
				Tags: []string{},
			},
		},
		{
			name: "given payload and location should return a domain capture with location without tags",
			payl: adding.Capture{
				Payload: adding.Payload{
					Payload: []domain.Metric{
						{Name: "power", Value: 10.0},
					},
				},
				Location: &adding.GeoLocation{LAT: f2P(1), LNG: f2P(1)},
			},
			expected: domain.Capture{
				Payload: domain.Payload{
					domain.Metric{Name: "power", Value: 10.0},
				},
				Location: &domain.Point{LAT: f2P(1), LNG: f2P(1)},
				Tags:     []string{},
			},
		},
		{
			name: "given payload and location should return a domain capture with location and tags",
			payl: adding.Capture{
				Payload: adding.Payload{
					Payload: []domain.Metric{
						{Name: "power", Value: 10.0},
					},
				},
				Location: &adding.GeoLocation{LAT: f2P(1), LNG: f2P(1), Elevation: f2P(1)},
				Tags:     []string{"at night"},
			},
			expected: domain.Capture{
				Payload: domain.Payload{
					domain.Metric{Name: "power", Value: 10.0},
				},
				Location: &domain.Point{LAT: f2P(1), LNG: f2P(1), Elevation: f2P(1)},
				Tags:     []string{"at night"},
			},
		},
	}

	repo := &domain.Repository{ID: kallax.NewULID()}
	s := adding.NewService(&mockStore{})

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			crrTime := time.Now()
			capt, err := s.AddCapture(repo, tc.payl)
			assert.Nil(t, err)

			assert.NotNil(t, capt.ID)
			assert.Equal(t, tc.expected.Payload, capt.Payload)
			assert.Equal(t, tc.expected.Location, capt.Location)
			assert.Equal(t, tc.expected.Tags, capt.Tags)
			assert.NotNil(t, capt.CreatedAt)
			assert.NotNil(t, capt.UpdatedAt)
			assert.Nil(t, capt.DeletedAt)

			assert.True(t, capt.Timestamp.After(crrTime))
		})
	}
}

func s2t(date string) time.Time {
	v, _ := time.Parse(time.RFC3339, date)
	return v
}

func TestServiceAddCaptureOKWithTimestamp(t *testing.T) {
	t.Parallel()

	expectedTime := s2t("1989-12-26T06:01:00.00Z")

	tt := adding.Timestamp{Date: s2n("1989-12-26T06:01:00.00Z")}
	adding.SetPostTimestampInstance(&tt, expectedTime)
	payl := adding.Capture{
		Payload: adding.Payload{
			Payload: []domain.Metric{
				{Name: "power", Value: 10.0},
			},
		},
		Timestamp: tt,
		Location:  &adding.GeoLocation{LAT: f2P(1), LNG: f2P(1), Elevation: f2P(1)},
		Tags:      []string{"at night"},
	}
	expected := domain.Capture{
		Payload: domain.Payload{
			domain.Metric{Name: "power", Value: 10.0},
		},
		Timestamp: expectedTime,
		Location:  &domain.Point{LAT: f2P(1), LNG: f2P(1), Elevation: f2P(1)},
		Tags:      []string{"at night"},
	}

	repo := &domain.Repository{ID: kallax.NewULID()}
	s := adding.NewService(&mockStore{})

	capt, err := s.AddCapture(repo, payl)
	assert.Nil(t, err)

	assert.NotNil(t, capt.ID)
	assert.Equal(t, expected.Payload, capt.Payload)
	assert.Equal(t, expected.Location, capt.Location)
	assert.Equal(t, expected.Timestamp, capt.Timestamp)
	assert.Equal(t, expected.Tags, capt.Tags)
	assert.NotNil(t, capt.CreatedAt)
	assert.NotNil(t, capt.UpdatedAt)
	assert.Nil(t, capt.DeletedAt)
}

func TestServiceAddCaptureErrWhenSaving(t *testing.T) {
	t.Parallel()
	s := adding.NewService(&mockStore{err: errors.New("test")})

	repo := &domain.Repository{ID: kallax.NewULID()}
	payl := adding.Capture{
		Payload: adding.Payload{
			Payload: []domain.Metric{
				{Name: "power", Value: 10.0},
			},
		},
	}

	_, err := s.AddCapture(repo, payl)
	assert.EqualError(t, err, "could not add capture: test")
}
