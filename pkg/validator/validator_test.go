package validator_test

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ifreddyrondon/capture/pkg/validator"
)

type mockValidator struct{ err error }

func (m mockValidator) OK() error { return m.err }

func TestStringValidatorDecode(t *testing.T) {
	t.Parallel()
	req := &http.Request{
		Method: "POST",
		Body:   ioutil.NopCloser(strings.NewReader(`{"msg":"hello"}`)),
	}
	err := validator.DefaultJSONValidator.Decode(req, new(mockValidator))
	assert.Nil(t, err)
}

func TestStringValidatorDecodeErr(t *testing.T) {
	t.Parallel()
	req := &http.Request{
		Method: "POST",
		Body:   ioutil.NopCloser(strings.NewReader("{")),
	}
	err := validator.DefaultJSONValidator.Decode(req, new(mockValidator))
	assert.EqualError(t, err, "cannot unmarshal json body")
}
