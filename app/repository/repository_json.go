package repository

import (
	"github.com/mailru/easyjson/jlexer"
	"github.com/markbates/validate"
)

const errNameRequired = "name must not be blank"

type repositoryJSON struct {
	Name   string `json:"name"`
	Shared bool   `json:"shared"`
}

func (r *repositoryJSON) IsValid(errors *validate.Errors) {
	if r.Name == "" {
		errors.Add("name", errNameRequired)
	}
}

func (r *repositoryJSON) toRepository() Repository {
	return Repository{Name: r.Name, Shared: r.Shared}
}

// UnmarshalJSON supports json.Unmarshaler interface
func (r *repositoryJSON) UnmarshalJSON(data []byte) error {
	l := jlexer.Lexer{Data: data}
	easyjsonD0cf849fDecodeGithubComIfreddyrondonCaptureRepository(&l, r)
	return l.Error()
}
