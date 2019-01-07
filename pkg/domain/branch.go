package domain

import kallax "gopkg.in/src-d/go-kallax.v1"

// Branch is a partial or full collection of captures within a repository.
type Branch struct {
	ID       kallax.ULID `json:"id"`
	Name     string      `json:"name"`
	Captures []Capture   `json:"captures"`
}
