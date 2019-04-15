package authorizing_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/src-d/go-kallax.v1"

	"github.com/ifreddyrondon/capture/pkg/authorizing"
	"github.com/ifreddyrondon/capture/pkg/domain"
)

func TestRepoPermission(t *testing.T) {
	sameID := kallax.NewULID()

	tt := []struct {
		name                    string
		repo                    *domain.Repository
		userID                  kallax.ULID
		expectedIsOwner         bool
		expectedIsOwnerOrPublic bool
	}{
		{
			name:                    "when repo dont belong to user and it's not public",
			repo:                    &domain.Repository{UserID: kallax.NewULID(), Visibility: domain.Private},
			userID:                  kallax.NewULID(),
			expectedIsOwner:         false,
			expectedIsOwnerOrPublic: false,
		},
		{
			name:                    "when repo dont belong to user but it's public",
			repo:                    &domain.Repository{UserID: kallax.NewULID(), Visibility: domain.Public},
			userID:                  kallax.NewULID(),
			expectedIsOwner:         false,
			expectedIsOwnerOrPublic: true,
		},
		{
			name:                    "when repo belong to user and it's not public",
			repo:                    &domain.Repository{UserID: sameID, Visibility: domain.Private},
			userID:                  sameID,
			expectedIsOwner:         true,
			expectedIsOwnerOrPublic: true,
		},
		{
			name:                    "when repo belong to user and it's not public",
			repo:                    &domain.Repository{UserID: sameID, Visibility: domain.Public},
			userID:                  sameID,
			expectedIsOwner:         true,
			expectedIsOwnerOrPublic: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			p := authorizing.RepoPermission(*tc.repo)
			assert.Equal(t, tc.expectedIsOwnerOrPublic, p.IsOwnerOrPublic(tc.userID))
			assert.Equal(t, tc.expectedIsOwner, p.IsOwner(tc.userID))
		})
	}
}
