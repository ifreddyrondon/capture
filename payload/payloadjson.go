package payload

import (
	jlexer "github.com/mailru/easyjson/jlexer"
)

type jsonPayload struct {
	Cap      map[string]interface{} `json:"cap"`
	Captures map[string]interface{} `json:"captures"`
	Data     map[string]interface{} `json:"data"`
	Payload  map[string]interface{} `json:"payload"`
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *jsonPayload) unmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson6ad23cceDecodeGithubComIfreddyrondonGocapturePayload(&r, v)
	return r.Error()
}
