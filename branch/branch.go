package branch

import (
	"github.com/markbates/going/defaults"

	"github.com/ifreddyrondon/gocapture/capture"
)

const defaultBranchName = "master"

// Branch represent a collection of captures.
type Branch struct {
	ID       string             `json:"id"`
	Name     string             `json:"name"`
	Captures []*capture.Capture `json:"captures"`
}

// New returns a new branch
func New(name string, captures ...*capture.Capture) *Branch {
	b := Branch{
		Name:     defaults.String(name, defaultBranchName),
		Captures: captures,
	}
	return &b
}
