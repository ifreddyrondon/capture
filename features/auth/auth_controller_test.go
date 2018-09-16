package auth_test

import (
	"net/http"
	"sync"
	"testing"

	"github.com/ifreddyrondon/capture/features/auth"
	"github.com/ifreddyrondon/capture/features/auth/authentication"
	"github.com/ifreddyrondon/capture/features/auth/authentication/strategy/basic"
	"github.com/ifreddyrondon/capture/features/auth/jwt"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/require"

	"github.com/ifreddyrondon/capture/features/user"

	"github.com/ifreddyrondon/bastion"
)

var (
	once sync.Once
	db   *gorm.DB
)

const (
	testUserEmail    = "test@example.com"
	testUserPassword = "b4KeHAYy3u9v=ZQX"
)

func getDB(t *testing.T) *gorm.DB {
	once.Do(func() {
		var err error
		db, err = gorm.Open("postgres", "postgres://localhost/captures_app_test?sslmode=disable")
		if err != nil {
			t.Fatal(err)
		}
	})
	return db
}

func setup(t *testing.T) (*bastion.Bastion, func()) {
	userStore := user.NewPGStore(getDB(t).Table("auth-users"))
	userStore.Migrate()
	teardown := func() { userStore.Drop() }
	userService := user.NewService(userStore)

	// save a user to test
	u := user.User{Email: testUserEmail}
	err := u.SetPassword(testUserPassword)
	require.Nil(t, err)
	userService.Save(&u)
	strategy := basic.New(userService)
	jwtService := jwt.NewService([]byte("test"), jwt.DefaultJWTExpirationDelta)
	middleware := authentication.Authenticate(strategy)
	controller := auth.NewController(middleware, jwtService)

	app := bastion.New()
	app.APIRouter.Mount("/auth/", controller.Router())

	return app, teardown
}

func TestBasicAuthentication(t *testing.T) {
	app, teardown := setup(t)
	defer teardown()

	payload := map[string]interface{}{"email": testUserEmail, "password": testUserPassword}
	e := bastion.Tester(t, app)
	e.POST("/auth/token-auth").WithJSON(payload).
		Expect().
		Status(http.StatusOK).
		JSON().Object().
		ContainsKey("token")
}
