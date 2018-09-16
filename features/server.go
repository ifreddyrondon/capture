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

	userService := cfg.Container.Get("user-service").(user.Service)
	authenticationStrategy := basic.New(userService)
	authenticationMiddleware := authentication.NewAuthentication(authenticationStrategy)
	jwtService := jwt.NewService([]byte("test"), jwt.DefaultJWTExpirationDelta)
	authController := auth.NewController(authenticationMiddleware, jwtService)
	authorizationMiddleware := authorization.NewAuthorization(jwtService)
	repoService := cfg.Container.Get("repo-service").(repository.Service)
	captureService := cfg.Container.Get("capture-service").(capture.Service)

	app.APIRouter.Mount("/users/", user.Routes(userService))
	app.APIRouter.Mount("/auth/", authController.Router())
	app.APIRouter.Mount("/captures/", capture.Routes(captureService))
	app.APIRouter.Mount("/branches/", branch.Routes())
	app.APIRouter.Mount("/repository/", repository.Routes(repoService, authorizationMiddleware.IsAuthorizedREQ, user.LoggedUser(userService)))
	return app
}
