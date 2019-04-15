package authorizing

import (
	"gopkg.in/src-d/go-kallax.v1"

	"github.com/ifreddyrondon/capture/pkg/domain"
)

type Permission interface {
	IsOwner(kallax.ULID) bool
	IsOwnerOrPublic(kallax.ULID) bool
}

type RepoPermission domain.Repository

func (c RepoPermission) IsOwner(ownerID kallax.ULID) bool {
	return ownerID == c.UserID
}

func (c RepoPermission) IsOwnerOrPublic(ownerID kallax.ULID) bool {
	return ownerID == c.UserID || c.Visibility == domain.Public
}
