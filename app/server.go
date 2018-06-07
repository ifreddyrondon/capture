package app

import (
	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/render/json"
	"github.com/ifreddyrondon/capture/app/auth/authorization"

	"github.com/ifreddyrondon/capture/app/auth"
	"github.com/ifreddyrondon/capture/app/auth/authentication"
	"github.com/ifreddyrondon/capture/app/auth/authentication/strategy/basic"
	"github.com/ifreddyrondon/capture/app/auth/jwt"
	"github.com/ifreddyrondon/capture/app/branch"
	"github.com/ifreddyrondon/capture/app/capture"
	"github.com/ifreddyrondon/capture/app/repository"
	"github.com/ifreddyrondon/capture/app/user"

	"github.com/jinzhu/gorm"
)

// New returns a bastion ready instance with all the app config
func New(db *gorm.DB) *bastion.Bastion {
	app := bastion.New(bastion.Options{})

	userStore := user.NewPGStore(db)
	userStore.Drop()
	userStore.Migrate()
	userService := user.NewService(userStore)
	userController := user.NewController(userService, json.NewRender)
	app.APIRouter.Mount("/users/", userController.Router())

	authenticationStrategy := basic.New(userService)
	authenticationMiddleware := authentication.NewAuthentication(authenticationStrategy, json.NewRender)

	jwtService := jwt.NewService([]byte("test"), jwt.DefaultJWTExpirationDelta)

	authController := auth.NewController(authenticationMiddleware, jwtService, json.NewRender)
	app.APIRouter.Mount("/auth/", authController.Router())

	authorizationMiddleware := authorization.NewAuthorization(jwtService, json.NewRender)

	repoStore := repository.NewPGStore(db)
	repoStore.Drop()
	repoStore.Migrate()
	repoService := repository.NewService(repoStore)
	repoController := repository.NewController(repoService, json.NewRender, authorizationMiddleware)
	app.APIRouter.Mount("/repository/", repoController.Router())

	captureStore := capture.NewPGStore(db)
	captureStore.Drop()
	captureStore.Migrate()
	captureService := capture.NewService(captureStore)
	captureController := capture.NewController(captureService, json.NewRender)
	app.APIRouter.Mount("/captures/", captureController.Router())

	branchController := branch.NewController(json.NewRender)
	app.APIRouter.Mount("/branches/", branchController.Router())
	return app
}
