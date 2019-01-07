package authorizing

import (
	"github.com/ifreddyrondon/capture/pkg/domain"
	"gopkg.in/src-d/go-kallax.v1"
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
