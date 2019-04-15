package adding_test

import (
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"gopkg.in/src-d/go-kallax.v1"

	"github.com/ifreddyrondon/capture/pkg/adding"
	"github.com/ifreddyrondon/capture/pkg/domain"
	"github.com/ifreddyrondon/capture/pkg/validator"
)

type mockMultiCaptureStore struct {
	err error
}

func (m *mockMultiCaptureStore) CreateCaptures(...domain.Capture) error { return m.err }

func TestServiceAddMultiCaptureOK(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name        string
		payl        adding.MultiCapture
		expectedLen int
		expected    []domain.Capture
	}{
		{
			name: "given a only name payload should return a domain capture without location and tags",
			payl: adding.MultiCapture{
				CapturesOK: []adding.Capture{
					{
						Payload: validator.Payload{
							Payload: []domain.Metric{{Name: "power", Value: 10.0}},
						},
					},
					{
						Payload: validator.Payload{
							Payload: []domain.Metric{{Name: "power", Value: 30.0}},
						},
					},
				},
			},
			expectedLen: 2,
			expected: []domain.Capture{
				{
					Payload: domain.Payload{
						domain.Metric{Name: "power", Value: 10.0},
					},
					Tags: []string{},
				},
				{
					Payload: domain.Payload{
						domain.Metric{Name: "power", Value: 30.0},
					},
					Tags: []string{},
				},
			},
		},
	}

	repo := &domain.Repository{ID: kallax.NewULID()}
	s := adding.NewMultiCaptureService(&mockMultiCaptureStore{})

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			crrTime := time.Now()
			captures, err := s.AddCaptures(repo, tc.payl)
			assert.Nil(t, err)
			assert.Len(t, captures, tc.expectedLen)

			for i, capt := range captures {
				assert.NotNil(t, capt.ID)
				assert.Equal(t, tc.expected[i].Payload, capt.Payload)
				assert.Equal(t, tc.expected[i].Location, capt.Location)
				assert.Equal(t, tc.expected[i].Tags, capt.Tags)
				assert.NotNil(t, capt.CreatedAt)
				assert.NotNil(t, capt.UpdatedAt)
				assert.Nil(t, capt.DeletedAt)
				assert.True(t, capt.Timestamp.After(crrTime))
			}
		})
	}
}

func TestServiceAddMultiCaptureErrWhenSaving(t *testing.T) {
	t.Parallel()
	s := adding.NewMultiCaptureService(&mockMultiCaptureStore{err: errors.New("test")})

	repo := &domain.Repository{ID: kallax.NewULID()}
	payl := adding.MultiCapture{
		CapturesOK: []adding.Capture{
			{
				Payload: validator.Payload{
					Payload: []domain.Metric{{Name: "power", Value: 10.0}},
				},
			},
		},
	}

	_, err := s.AddCaptures(repo, payl)
	assert.EqualError(t, err, "could not add captures: test")
}
