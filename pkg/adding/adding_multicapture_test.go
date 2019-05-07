package adding_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/ifreddyrondon/bastion/binder"
	"github.com/stretchr/testify/assert"

	"github.com/ifreddyrondon/capture/pkg/adding"
	"github.com/ifreddyrondon/capture/pkg/domain"
	"github.com/ifreddyrondon/capture/pkg/validator"
)

func TestValidateMultiCaptureOK(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name       string
		body       string
		expectedOK int
		expected   adding.MultiCapture
	}{
		{
			name: "decode multi capture with no ignore errors and all valid",
			body: `{
					"captures":[
						{"payload":[{"name":"power","value":10}]},
						{"payload":[{"name":"power","value":30}]}
					],
					"ignore_errors":false
				}`,
			expectedOK: 2,
			expected: adding.MultiCapture{
				IgnoreErrors: false,
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
		},
		{
			name: "decode multi capture with ignore errors and one invalid",
			body: `{
					"captures":[
						{"payload":[{"name":"power","value":10}]},
						{"payload":[]}
					],
					"ignore_errors":true
				}`,
			expectedOK: 1,
			expected: adding.MultiCapture{
				IgnoreErrors: true,
				CapturesOK: []adding.Capture{
					{
						Payload: validator.Payload{
							Payload: []domain.Metric{{Name: "power", Value: 10.0}},
						},
					},
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			r, _ := http.NewRequest("POST", "/", strings.NewReader(tc.body))

			var multiCapture adding.MultiCapture
			err := binder.JSON.FromReq(r, &multiCapture)
			assert.Nil(t, err)
			assert.Equal(t, tc.expectedOK, len(multiCapture.CapturesOK))
			for i, c := range multiCapture.CapturesOK {
				assert.Equal(t, tc.expected.CapturesOK[i].Payload, c.Payload)
				assert.Equal(t, tc.expected.CapturesOK[i].Location, c.Location)
				assert.Equal(t, tc.expected.CapturesOK[i].Timestamp.Timestamp, c.Timestamp.Timestamp)
				assert.Equal(t, tc.expected.CapturesOK[i].Timestamp.Date, c.Timestamp.Date)
				assert.Equal(t, tc.expected.CapturesOK[i].Tags, c.Tags)
			}
		})
	}
}

func TestValidationMultiCaptureFails(t *testing.T) {
	t.Parallel()
	tt := []struct {
		name string
		body string
		errs []string
	}{
		{
			name: "decode multi capture when missing payload",
			body: `{
					"ignore_errors":false
				}`,
			errs: []string{"captures value must not be blank or empty"},
		},
		{
			name: "decode multi capture when empty payload",
			body: `{
					"captures":[],
					"ignore_errors":false
				}`,
			errs: []string{"captures value must not be blank or empty"},
		},
		{
			name: "decode multi capture when max allowed",
			body: `{
					"captures":[{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{},{}],
					"ignore_errors":false
				}`,
			errs: []string{"the maximum amount of allowed captures is 50"},
		},
		{
			name: "decode multi capture with ignore errors false and one invalid",
			body: `{
					"captures":[
						{"payload":[{"name":"power","value":10}]},
						{"payload":[]}
					],
					"ignore_errors":false
				}`,
			errs: []string{"capture 1: payload value must not be blank"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			r, _ := http.NewRequest("POST", "/", strings.NewReader(tc.body))
			var multiCapture adding.MultiCapture
			err := binder.JSON.FromReq(r, &multiCapture)
			for _, e := range tc.errs {
				assert.Contains(t, err.Error(), e)
			}
		})
	}
}
