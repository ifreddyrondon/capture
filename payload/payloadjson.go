package payload

import "github.com/mailru/easyjson/jlexer"

type jsonPayload struct {
	Cap      []*Metric `json:"cap"`
	Captures []*Metric `json:"captures"`
	Data     []*Metric `json:"data"`
	Payload  []*Metric `json:"payload"`
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *jsonPayload) unmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson6ad23cceDecodeGithubComIfreddyrondonGocapturePayload(&r, v)
	return r.Error()
}
