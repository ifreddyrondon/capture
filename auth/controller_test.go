package auth_test

import (
	"net/http"
	"sync"
	"testing"

	"github.com/ifreddyrondon/gocapture/jwt"

	"github.com/ifreddyrondon/gocapture/app"
	"github.com/ifreddyrondon/gocapture/auth"
	"github.com/ifreddyrondon/gocapture/auth/strategy/basic"
	"github.com/ifreddyrondon/gocapture/database"
	"github.com/ifreddyrondon/gocapture/user"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/require"

	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/render/json"
)

var (
	once sync.Once
	db   *gorm.DB
)

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

	strategy := basic.Strategy{
		Render:        json.NewRender,
		UserKey:       app.ContextKey("user"),
		GetterService: userService,
	}

	jwtService := jwt.NewService([]byte("test"), jwt.DefaultJWTExpirationDelta, json.NewRender)

	controller := auth.Controller{
		Strategy: strategy,
		Render:   json.NewRender,
		UserKey:  app.ContextKey("user"),
		JWT:      jwtService,
	}

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
