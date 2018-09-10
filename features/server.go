package features

import (
	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/capture/config"
	"github.com/ifreddyrondon/capture/features/auth"
	"github.com/ifreddyrondon/capture/features/auth/authentication"
	"github.com/ifreddyrondon/capture/features/auth/authentication/strategy/basic"
	"github.com/ifreddyrondon/capture/features/auth/authorization"
	"github.com/ifreddyrondon/capture/features/auth/jwt"
	"github.com/ifreddyrondon/capture/features/branch"
	"github.com/ifreddyrondon/capture/features/capture"
	"github.com/ifreddyrondon/capture/features/repository"
	"github.com/ifreddyrondon/capture/features/user"
)

// New returns a bastion ready instance with all the features config
func New(cfg *config.Config) *bastion.Bastion {
	app := bastion.New(bastion.Addr(cfg.ADDR))

	userService := cfg.Container.Get("user-service").(user.Store)
	app.APIRouter.Mount("/users/", user.Routes(userService))

	authenticationStrategy := basic.New(userService)
	authenticationMiddleware := authentication.NewAuthentication(authenticationStrategy)

	jwtService := jwt.NewService([]byte("test"), jwt.DefaultJWTExpirationDelta)

	authController := auth.NewController(authenticationMiddleware, jwtService)
	app.APIRouter.Mount("/auth/", authController.Router())

	authorizationMiddleware := authorization.NewAuthorization(jwtService)

	repoService := cfg.Container.Get("repo-service").(repository.Store)
	repoController := repository.NewController(repoService, authorizationMiddleware, userService)
	app.APIRouter.Mount("/repository/", repoController.Router())

	captureService := cfg.Container.Get("capture-service").(capture.Store)
	captureController := capture.NewController(captureService)
	app.APIRouter.Mount("/captures/", captureController.Router())

	branchController := branch.NewController()
	app.APIRouter.Mount("/branches/", branchController.Router())
	return app
}
