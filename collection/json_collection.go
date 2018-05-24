package collection

import (
	"github.com/mailru/easyjson/jlexer"
	"github.com/markbates/validate"
)

const errNameRequired = "name must not be blank"

type jsonCollection struct {
	Name string `json:"name"`
}

func (v *jsonCollection) IsValid(errors *validate.Errors) {
	if v.Name == "" {
		errors.Add("name", errNameRequired)
	}
}

func (v *jsonCollection) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD0cf849fDecodeGithubComIfreddyrondonGocaptureCollection(&r, v)
	return r.Error()
}
