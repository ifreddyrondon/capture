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
	"github.com/ifreddyrondon/capture/features/repository"
	"github.com/ifreddyrondon/capture/features/user"
)

func router(cfg *config.Config) http.Handler {
	r := chi.NewRouter()

	userService := cfg.Container.Get("user-service").(user.Service)
	repoService := cfg.Container.Get("repo-service").(repository.Service)

	jwtService := jwt.NewService([]byte("test"), jwt.DefaultJWTExpirationDelta)
	authenticationStrategy := basic.New(userService)
	authRoutes := auth.Routes(authentication.Authenticate(authenticationStrategy), jwtService)

	userRoutes := cfg.Container.Get("user-routes").(http.Handler)
	captureRoutes := cfg.Container.Get("capture-routes").(http.Handler)
	branchRoutes := cfg.Container.Get("branch-routes").(http.Handler)

	r.Mount("/users/", userRoutes)
	r.Mount("/auth/", authRoutes)
	r.Mount("/captures/", captureRoutes)
	r.Mount("/branches/", branchRoutes)
	r.Mount("/repository/", repository.Routes(repoService, authorization.IsAuthorizedREQ(jwtService), user.LoggedUser(userService)))
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
