package auth_test

import (
	"net/http"
	"sync"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/require"

	"github.com/ifreddyrondon/capture/app/auth"
	"github.com/ifreddyrondon/capture/app/auth/strategy/basic"
	"github.com/ifreddyrondon/capture/app/jwt"
	"github.com/ifreddyrondon/capture/app/user"
	"github.com/ifreddyrondon/capture/database"

	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/render/json"
)

var (
	once sync.Once
	db   *gorm.DB
)

const (
	testUserEmail    = "test@example.com"
	testUserPassword = "b4KeHAYy3u9v=ZQX"
)

func getDB() *gorm.DB {
	once.Do(func() {
		ds := database.Open("postgres://localhost/captures_app_test?sslmode=disable")
		db = ds.DB
	})
	return db
}

func setup(t *testing.T) (*bastion.Bastion, func()) {
	userStore := user.NewPGStore(getDB().Table("auth-users"))
	userStore.Migrate()
	teardown := func() { userStore.Drop() }
	userService := user.NewService(userStore)

	// save a user to test
	u := user.User{Email: testUserEmail}
	err := u.SetPassword(testUserPassword)
	require.Nil(t, err)
	userService.Save(&u)
	strategy := basic.NewStrategy(json.NewRender, userService)
	jwtService := jwt.NewService([]byte("test"), jwt.DefaultJWTExpirationDelta, json.NewRender)
	controller := auth.NewController(strategy, jwtService, json.NewRender)

	app := bastion.New(bastion.Options{})
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
