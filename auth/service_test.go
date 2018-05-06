package auth_test

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"

	"github.com/ifreddyrondon/gocapture/auth"
	"github.com/ifreddyrondon/gocapture/database"

	"github.com/ifreddyrondon/gocapture/user"
	"github.com/jinzhu/gorm"
)

var once sync.Once
var db *gorm.DB

const (
	testUserEmail    = "test@test.com"
	testUserPassword = "b4KeHAYy3u9v=ZQX"
)

func getDB() *gorm.DB {
	once.Do(func() {
		ds := database.Open("postgres://localhost/captures_app_test?sslmode=disable")
		db = ds.DB
	})
	return db
}

func setupService(t *testing.T) (*auth.PGAuthService, func()) {
	userService := user.PGService{DB: getDB().Table("auth_users")}
	userService.Migrate()

	// save a user to test
	u := user.User{Email: testUserEmail}
	err := u.SetPassword(testUserPassword)
	require.Nil(t, err)
	userService.Save(&u)

	teardown := func() { userService.Drop() }
	service := auth.PGAuthService{GetterService: &userService}

	return &service, teardown
}

func TestAuthenticateSuccess(t *testing.T) {
	service, teardown := setupService(t)
	defer teardown()

	credentials := auth.BasicAuthCrendential{Email: testUserEmail, Password: testUserPassword}
	u, err := service.Authenticate(&credentials)
	assert.Nil(t, err)
	assert.NotNil(t, u.ID)
	assert.Equal(t, u.Email, testUserEmail)
}

func TestAuthenticateFailure(t *testing.T) {
	service, teardown := setupService(t)
	defer teardown()

	tt := []struct {
		name        string
		credentials *auth.BasicAuthCrendential
		err         string
	}{
		{
			"email not found",
			&auth.BasicAuthCrendential{Email: "invalid@email.com", Password: "123"},
			"user not found",
		},
		{
			"invalid password",
			&auth.BasicAuthCrendential{Email: testUserEmail, Password: "123"},
			"invalid email or password",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			_, err := service.Authenticate(tc.credentials)
			assert.Error(t, err)
			assert.Equal(t, tc.err, err.Error())
		})
	}
}
