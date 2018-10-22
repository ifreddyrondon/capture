package decoder

import (
	"strings"
	"time"

	"github.com/ifreddyrondon/capture/features"
	"github.com/markbates/validate"
	"gopkg.in/src-d/go-kallax.v1"
)

const (
	defaultCrrBranchFieldValue = "master"
	errNameRequired            = "name must not be blank"
)

type PostRepository struct {
	Name   *string `json:"name"`
	Shared *bool   `json:"shared"`
}

func (r PostRepository) OK() error {
	e := validate.NewErrors()
	if r.Name == nil {
		e.Add("name", errNameRequired)
	} else if len(strings.TrimSpace(*r.Name)) == 0 {
		e.Add("name", errNameRequired)
	}
	if e.HasAny() {
		return e
	}

	return nil
}

func (r PostRepository) GetRepository() features.Repository {
	now := time.Now()
	repo := features.Repository{
		ID:            kallax.NewULID(),
		Name:          *r.Name,
		CurrentBranch: defaultCrrBranchFieldValue,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	if r.Shared == nil {
		repo.Shared = true
	} else {
		repo.Shared = *r.Shared
	}

	return repo
}
