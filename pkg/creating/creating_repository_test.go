package creating_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ifreddyrondon/capture/pkg/creating"
)

func TestValidatePayloadOK(t *testing.T) {
	t.Parallel()
	name, public, private := "test_repository", "public", "private"

	tt := []struct {
		name     string
		body     string
		expected creating.Payload
	}{
		{
			name:     "decode repo without shared",
			body:     `{"name":"test_repository"}`,
			expected: creating.Payload{Name: &name, Visibility: nil},
		},
		{
			name:     "decode repo with shared true",
			body:     `{"name":"test_repository","visibility":"public"}`,
			expected: creating.Payload{Name: &name, Visibility: &public},
		},
		{
			name:     "decode repo with shared true",
			body:     `{"name":"test_repository","visibility":"private"}`,
			expected: creating.Payload{Name: &name, Visibility: &private},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			r, _ := http.NewRequest("POST", "/", strings.NewReader(tc.body))

			var p creating.Payload
			err := creating.Validator.Decode(r, &p)
			assert.Nil(t, err)
			assert.Equal(t, tc.expected.Name, p.Name)
			assert.Equal(t, tc.expected.Visibility, p.Visibility)
		})
	}
}

func TestValidatePayloadFails(t *testing.T) {
	t.Parallel()
	tt := []struct {
		name string
		body string
		err  string
	}{
		{
			name: "decode payload name missing",
			body: `{}`,
			err:  "name must not be blank",
		},
		{
			name: "decode payload name empty",
			body: `{"name":""}`,
			err:  "name must not be blank",
		},
		{
			name: "decode payload name empty v2",
			body: `{"name":"   "}`,
			err:  "name must not be blank",
		},
		{
			name: "decode payload with not allowed visibility",
			body: `{"name":"foo","visibility":"protected"}`,
			err:  "not allowed visibility type. it Could be one of public, or private. Default: public",
		},
		{
			name: "decode payload with not allowed visibility because empty",
			body: `{"name":"foo","visibility":""}`,
			err:  "not allowed visibility type. it Could be one of public, or private. Default: public",
		},
		{
			name: "invalid payload payload",
			body: `.`,
			err:  "cannot unmarshal json into valid repository",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			r, _ := http.NewRequest("POST", "/", strings.NewReader(tc.body))

			var p creating.Payload
			err := creating.Validator.Decode(r, &p)
			assert.EqualError(t, err, tc.err)
		})
	}
}
