package app

import (
	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/capture/app/auth"
	"github.com/ifreddyrondon/capture/app/auth/authentication"
	"github.com/ifreddyrondon/capture/app/auth/authentication/strategy/basic"
	"github.com/ifreddyrondon/capture/app/auth/authorization"
	"github.com/ifreddyrondon/capture/app/auth/jwt"
	"github.com/ifreddyrondon/capture/app/branch"
	"github.com/ifreddyrondon/capture/app/capture"
	"github.com/ifreddyrondon/capture/app/repository"
	"github.com/ifreddyrondon/capture/app/user"
	"github.com/ifreddyrondon/capture/internal/config"
)

// New returns a bastion ready instance with all the app config
func New(cfg *config.Config) *bastion.Bastion {
	app := bastion.New(bastion.Addr(cfg.ADDR))

	userStore := user.NewPGStore(cfg.Database)
	userStore.Drop()
	userStore.Migrate()
	userService := user.NewService(userStore)
	userController := user.NewController(userService)
	userMiddleware := user.NewMiddleware(userService)
	app.APIRouter.Mount("/users/", userController.Router())

	authenticationStrategy := basic.New(userService)
	authenticationMiddleware := authentication.NewAuthentication(authenticationStrategy)

	jwtService := jwt.NewService([]byte("test"), jwt.DefaultJWTExpirationDelta)

	authController := auth.NewController(authenticationMiddleware, jwtService)
	app.APIRouter.Mount("/auth/", authController.Router())

	authorizationMiddleware := authorization.NewAuthorization(jwtService)

	repoStore := repository.NewPGStore(cfg.Database)
	repoStore.Drop()
	repoStore.Migrate()
	repoService := repository.NewService(repoStore)
	repoController := repository.NewController(repoService, authorizationMiddleware, userMiddleware)
	app.APIRouter.Mount("/repository/", repoController.Router())

	captureStore := capture.NewPGStore(cfg.Database)
	captureStore.Drop()
	captureStore.Migrate()
	captureService := capture.NewService(captureStore)
	captureController := capture.NewController(captureService)
	app.APIRouter.Mount("/captures/", captureController.Router())

	branchController := branch.NewController()
	app.APIRouter.Mount("/branches/", branchController.Router())
	return app
}
