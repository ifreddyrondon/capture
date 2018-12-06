package config

import (
	"net/http"
	"time"

	"github.com/ifreddyrondon/capture/pkg/authenticating"
	"github.com/ifreddyrondon/capture/pkg/authorizing"
	"github.com/ifreddyrondon/capture/pkg/branch"
	"github.com/ifreddyrondon/capture/pkg/capture"
	"github.com/ifreddyrondon/capture/pkg/http/rest"
	"github.com/ifreddyrondon/capture/pkg/http/rest/middleware"
	"github.com/ifreddyrondon/capture/pkg/multipost"
	"github.com/ifreddyrondon/capture/pkg/repository"
	"github.com/ifreddyrondon/capture/pkg/storage/postgres"
	"github.com/ifreddyrondon/capture/pkg/token"
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
				middle := cfg.Resources.Get("authorizeReq-middleware").(func(next http.Handler) http.Handler)
				service := cfg.Resources.Get("repository-service").(repository.Service)
				return repository.UserRoutes(service, middle), nil
			},
		},
		{
			Name:  "repositories-routes",
			Scope: di.App,
			Build: func(ctn di.Container) (interface{}, error) {
				middle := cfg.Resources.Get("authorizeReq-middleware").(func(next http.Handler) http.Handler)
				service := cfg.Resources.Get("repository-service").(repository.Service)
				return repository.Routes(service, middle), nil
			},
		},
		{
			Name:  "multipost-routes",
			Scope: di.App,
			Build: func(ctn di.Container) (interface{}, error) {
				return multipost.Routes(), nil
			},
		},

		// resources for DDD migration
		{
			Name:  "postgres-storage",
			Scope: di.App,
			Build: func(ctn di.Container) (i interface{}, e error) {
				database := cfg.Resources.Get("database").(*gorm.DB)
				return postgres.NewPGStorage(database), nil
			},
		},
		{
			Name:  "jwt-service",
			Scope: di.App,
			Build: func(ctn di.Container) (i interface{}, e error) {
				duration := time.Duration(cfg.JWTExpirationDelta) * time.Second
				return token.NewJWTService(cfg.JWTSigningKey, duration), nil
			},
		},
		{
			Name:  "authenticating-routes",
			Scope: di.App,
			Build: func(ctn di.Container) (i interface{}, e error) {
				tokenService := cfg.Resources.Get("jwt-service").(authenticating.TokenService)
				store := cfg.Resources.Get("postgres-storage").(authenticating.Store)
				s := authenticating.NewService(tokenService, store)
				return rest.AuthenticatingRoutes(s), nil
			},
		},
		{
			Name:  "authorizeReq-middleware",
			Scope: di.App,
			Build: func(ctn di.Container) (i interface{}, e error) {
				tokenService := cfg.Resources.Get("jwt-service").(authorizing.TokenService)
				store := cfg.Resources.Get("postgres-storage").(authorizing.Store)
				s := authorizing.NewService(tokenService, store)
				return middleware.AuthorizeReq(s), nil
			},
		},
	}

	builder.Add(definitions...)
	return builder.Build()
}
