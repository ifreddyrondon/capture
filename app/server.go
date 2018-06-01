package app

import (
	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/render/json"
	"github.com/ifreddyrondon/gocapture/auth"
	"github.com/ifreddyrondon/gocapture/auth/strategy/basic"
	"github.com/ifreddyrondon/gocapture/branch"
	"github.com/ifreddyrondon/gocapture/capture"
	"github.com/ifreddyrondon/gocapture/jwt"
	"github.com/ifreddyrondon/gocapture/repository"
	"github.com/ifreddyrondon/gocapture/user"
	"github.com/jinzhu/gorm"
)

// ContextKey represent a key in a context.ContextValue
type ContextKey string

func (c ContextKey) String() string {
	return string(c)
}

// New returns a bastion ready instance with all the app config
func New(db *gorm.DB) *bastion.Bastion {
	app := bastion.New(bastion.Options{})

	userStore := user.NewPGStore(db)
	userStore.Drop()
	userStore.Migrate()
	userService := user.NewService(userStore)
	userController := user.NewController(userService, json.NewRender)
	app.APIRouter.Mount("/users/", userController.Router())

	strategy := basic.NewStrategy(json.NewRender, userService)
	jwtService := jwt.NewService([]byte("test"), jwt.DefaultJWTExpirationDelta, json.NewRender)
	authController := auth.NewController(strategy, jwtService, json.NewRender)
	app.APIRouter.Mount("/auth/", authController.Router())

	repoStore := repository.NewPGStore(db)
	repoStore.Drop()
	repoStore.Migrate()
	repoService := repository.NewService(repoStore)
	repoController := repository.NewController(repoService, json.NewRender)
	app.APIRouter.Mount("/repository/", repoController.Router())

	captureStore := capture.NewPGStore(db)
	captureStore.Drop()
	captureStore.Migrate()
	captureService := capture.NewService(captureStore)
	captureController := capture.NewController(captureService, json.NewRender, ContextKey("capture"))
	app.APIRouter.Mount("/captures/", captureController.Router())

	branchController := branch.NewController(json.NewRender)
	app.APIRouter.Mount("/branches/", branchController.Router())
	return app
}
