package numberlist

import "github.com/mailru/easyjson/jlexer"

type jsonPayload struct {
	Cap      []float64 `json:"cap"`
	Captures []float64 `json:"captures"`
	Data     []float64 `json:"data"`
	Payload  []float64 `json:"payload"`
}

func (v *jsonPayload) unmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC80ae7adDecodeGithubComIfreddyrondonGocapturePayload(&r, v)
	return r.Error()
}
