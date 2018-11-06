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
	errVisibilityNotAllowed    = "not allowed visibility type. it Could be one of public, or private. Default: public"
)

type PostRepository struct {
	Name       *string `json:"name"`
	Visibility *string `json:"visibility"`
}

func (r PostRepository) OK() error {
	e := validate.NewErrors()
	if r.Name == nil {
		e.Add("name", errNameRequired)
	} else if len(strings.TrimSpace(*r.Name)) == 0 {
		e.Add("name", errNameRequired)
	}

	if r.Visibility != nil {
		if !features.AllowedVisibility(*r.Visibility) {
			e.Add("name", errVisibilityNotAllowed)
		}
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
	if r.Visibility == nil {
		repo.Visibility = features.Public
	} else {
		repo.Visibility = features.Visibility(*r.Visibility)
	}

	return repo
}
