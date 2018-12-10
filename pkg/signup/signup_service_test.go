package signup_test

import (
	"testing"

	"github.com/ifreddyrondon/capture/pkg"
	"github.com/ifreddyrondon/capture/pkg/signup"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

type mockStore struct {
	usr *pkg.User
	err error
}

func (m *mockStore) SaveUser(user *pkg.User) error { return m.err }

func string2pointer(v string) *string { return &v }

func TestServiceEnrollUserOK(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name string
		payl signup.Payload
	}{
		{
			name: "should enroll user with just email",
			payl: signup.Payload{Email: string2pointer("ifreddyrondon@gmail.com")},
		},
		{
			name: "should enroll user with email and password",
			payl: signup.Payload{
				Email:    string2pointer("ifreddyrondon@gmail.com"),
				Password: string2pointer("1"),
			},
		},
	}

	s := signup.NewService(&mockStore{})
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			u, err := s.EnrollUser(tc.payl)
			assert.Nil(t, err)
			assert.Equal(t, "ifreddyrondon@gmail.com", u.Email)
			assert.NotNil(t, u.ID)
			assert.NotNil(t, u.CreatedAt)
			assert.NotNil(t, u.UpdatedAt)
		})
	}
}

type uniqueConstraintErr string

func (u uniqueConstraintErr) Error() string          { return string(u) }
func (u uniqueConstraintErr) UniqueConstraint() bool { return true }

type conflictErr interface {
	Conflict() bool
}

func TestServiceEnrollUserErrWhenDuplicatedEmail(t *testing.T) {
	t.Parallel()
	s := signup.NewService(&mockStore{err: uniqueConstraintErr("duplicated email")})
	payl := signup.Payload{Email: string2pointer("ifreddyrondon@gmail.com")}
	_, err := s.EnrollUser(payl)
	assert.EqualError(t, err, "email ifreddyrondon@gmail.com already exist")
	confErr, ok := errors.Cause(err).(conflictErr)
	assert.True(t, ok)
	assert.True(t, confErr.Conflict())
}

func TestServiceEnrollUserErrWhenSaving(t *testing.T) {
	t.Parallel()
	s := signup.NewService(&mockStore{err: errors.New("test")})
	payl := signup.Payload{Email: string2pointer("ifreddyrondon@gmail.com")}
	_, err := s.EnrollUser(payl)
	assert.EqualError(t, err, "could not save user: test")
}
