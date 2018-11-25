package config

import (
	"net/http"

	"github.com/ifreddyrondon/capture/pkg/auth"
	"github.com/ifreddyrondon/capture/pkg/auth/authentication"
	"github.com/ifreddyrondon/capture/pkg/auth/authentication/strategy/basic"
	"github.com/ifreddyrondon/capture/pkg/auth/authorization"
	"github.com/ifreddyrondon/capture/pkg/auth/jwt"
	"github.com/ifreddyrondon/capture/pkg/branch"
	"github.com/ifreddyrondon/capture/pkg/capture"
	"github.com/ifreddyrondon/capture/pkg/multipost"
	"github.com/ifreddyrondon/capture/pkg/repository"
	"github.com/ifreddyrondon/capture/pkg/user"
	"github.com/jinzhu/gorm"
	"github.com/sarulabs/di"
)

func getResources(cfg *Config) di.Container {
	builder, _ := di.NewBuilder()
	definitions := []di.Def{
		{
			Name:  "database",
			Scope: di.App,
			Build: func(ctn di.Container) (interface{}, error) {
				return gorm.Open("postgres", cfg.Constants.PG)
			},
			Close: func(obj interface{}) error {
				return obj.(*gorm.DB).Close()
			},
		},
		{
			Name:  "user-service",
			Scope: di.App,
			Build: func(ctn di.Container) (interface{}, error) {
				database := cfg.Resources.Get("database").(*gorm.DB)
				store := user.NewPGStore(database)
				store.Drop()
				store.Migrate()
				return user.NewService(store), nil
			},
		},
		{
			Name:  "jwt-service",
			Scope: di.App,
			Build: func(ctn di.Container) (interface{}, error) {
				return jwt.NewService([]byte("test"), jwt.DefaultJWTExpirationDelta), nil
			},
		},
		{
			Name:  "is-auth-middleware",
			Scope: di.App,
			Build: func(ctn di.Container) (interface{}, error) {
				jwtService := cfg.Resources.Get("jwt-service").(*jwt.Service)
				return authorization.IsAuthorizedREQ(jwtService), nil
			},
		},
		{
			Name:  "logged-user-middleware",
			Scope: di.App,
			Build: func(ctn di.Container) (interface{}, error) {
				userService := cfg.Resources.Get("user-service").(user.Service)
				return user.LoggedUser(userService), nil
			},
		},
		{
			Name:  "auth-routes",
			Scope: di.App,
			Build: func(ctn di.Container) (interface{}, error) {
				userService := cfg.Resources.Get("user-service").(user.Service)
				jwtService := cfg.Resources.Get("jwt-service").(*jwt.Service)
				strategy := basic.New(userService)

				return auth.Routes(authentication.Authenticate(strategy), jwtService), nil
			},
		},
		{
			Name:  "user-routes",
			Scope: di.App,
			Build: func(ctn di.Container) (interface{}, error) {
				userService := cfg.Resources.Get("user-service").(user.Service)
				return user.Routes(userService), nil
			},
		},
		{
			Name:  "capture-routes",
			Scope: di.App,
			Build: func(ctn di.Container) (interface{}, error) {
				database := cfg.Resources.Get("database").(*gorm.DB)
				store := capture.NewPGStore(database)
				store.Drop()
				store.Migrate()
				return capture.Routes(store), nil
			},
		},
		{
			Name:  "branch-routes",
			Scope: di.App,
			Build: func(ctn di.Container) (interface{}, error) {
				return branch.Routes(), nil
			},
		},
		{
			Name:  "repository-store",
			Scope: di.App,
			Build: func(ctn di.Container) (interface{}, error) {
				database := cfg.Resources.Get("database").(*gorm.DB)
				store := repository.NewPGStore(database)
				store.Drop()
				store.Migrate()
				return store, nil
			},
		},
		{
			Name:  "repository-service",
			Scope: di.App,
			Build: func(ctn di.Container) (interface{}, error) {
				store := cfg.Resources.Get("repository-store").(repository.Store)
				return repository.Service{Store: store}, nil
			},
		},
		{
			Name:  "user-repo-routes",
			Scope: di.App,
			Build: func(ctn di.Container) (interface{}, error) {
				loggedUser := cfg.Resources.Get("logged-user-middleware").(func(next http.Handler) http.Handler)
				isAuth := cfg.Resources.Get("is-auth-middleware").(func(next http.Handler) http.Handler)
				service := cfg.Resources.Get("repository-service").(repository.Service)
				return repository.UserRoutes(service, isAuth, loggedUser), nil
			},
		},
		{
			Name:  "repositories-routes",
			Scope: di.App,
			Build: func(ctn di.Container) (interface{}, error) {
				loggedUser := cfg.Resources.Get("logged-user-middleware").(func(next http.Handler) http.Handler)
				isAuth := cfg.Resources.Get("is-auth-middleware").(func(next http.Handler) http.Handler)
				service := cfg.Resources.Get("repository-service").(repository.Service)
				return repository.Routes(service, isAuth, loggedUser), nil
			},
		},
		{
			Name:  "multipost-routes",
			Scope: di.App,
			Build: func(ctn di.Container) (interface{}, error) {
				return multipost.Routes(), nil
			},
		},
	}

	builder.Add(definitions...)
	return builder.Build()
}
