package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
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

func router(cfg *config.Config) http.Handler {
	r := chi.NewRouter()

	userService := cfg.Container.Get("user-service").(user.Service)
	authenticationStrategy := basic.New(userService)
	authenticationMiddleware := authentication.NewAuthentication(authenticationStrategy)
	jwtService := jwt.NewService([]byte("test"), jwt.DefaultJWTExpirationDelta)
	authController := auth.NewController(authenticationMiddleware, jwtService)
	authorizationMiddleware := authorization.NewAuthorization(jwtService)
	repoService := cfg.Container.Get("repo-service").(repository.Service)
	captureService := cfg.Container.Get("capture-service").(capture.Service)

	r.Mount("/users/", user.Routes(userService))
	r.Mount("/auth/", authController.Router())
	r.Mount("/captures/", capture.Routes(captureService))
	r.Mount("/branches/", branch.Routes())
	r.Mount("/repository/", repository.Routes(repoService, authorizationMiddleware.IsAuthorizedREQ, user.LoggedUser(userService)))
	return r
}

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Panicln("Configuration error", err)
	}

	app := bastion.New(bastion.Addr(cfg.ADDR))
	app.APIRouter.Mount("/", router(cfg))
	app.RegisterOnShutdown(cfg.OnShutdown)
	if err := app.Serve(); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
}
