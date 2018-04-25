package capture

import "github.com/mailru/easyjson/jlexer"

type tagsJSON struct {
	Tags []string `json:"tags"`
}

func (v *tagsJSON) unmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson247433ddDecodeGithubComIfreddyrondonGocaptureCapture(&r, v)
	if v.Tags == nil {
		v.Tags = []string{}
	}
	return r.Error()
}
