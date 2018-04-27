package payload

import (
	"github.com/mailru/easyjson/jlexer"
	"github.com/mailru/easyjson/jwriter"
)

// Metric represent a value capture from a device/sensor
type Metric struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

// MarshalJSON supports json.Marshaler interface
func (v Metric) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson9478868cEncodeGithubComIfreddyrondonGocapturePayload(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Metric) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson9478868cDecodeGithubComIfreddyrondonGocapturePayload(&r, v)
	return r.Error()
}
