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

func (r PostRepository) ok() error {
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

func (r PostRepository) repository(repo *features.Repository) {
	repo.ID = kallax.NewULID()
	repo.Name = *r.Name
	if r.Shared == nil {
		repo.Shared = true
	} else {
		repo.Shared = *r.Shared
	}
	repo.CurrentBranch = defaultCrrBranchFieldValue
	now := time.Now()
	repo.CreatedAt = now
	repo.UpdatedAt = now
}
