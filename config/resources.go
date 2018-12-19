package config

import (
	"net/http"
	"time"

	"github.com/ifreddyrondon/capture/pkg/listing"
	"github.com/ifreddyrondon/capture/pkg/storage/postgres/repo"

	"github.com/ifreddyrondon/capture/pkg/creating"

	"github.com/ifreddyrondon/capture/pkg/authenticating"
	"github.com/ifreddyrondon/capture/pkg/authorizing"
	"github.com/ifreddyrondon/capture/pkg/branch"
	"github.com/ifreddyrondon/capture/pkg/capture"
	"github.com/ifreddyrondon/capture/pkg/http/rest"
	"github.com/ifreddyrondon/capture/pkg/http/rest/middleware"
	"github.com/ifreddyrondon/capture/pkg/multipost"
	"github.com/ifreddyrondon/capture/pkg/repository"
	"github.com/ifreddyrondon/capture/pkg/signup"
	"github.com/ifreddyrondon/capture/pkg/storage/postgres/user"
	"github.com/ifreddyrondon/capture/pkg/token"
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
			Name:  "repositories-routes",
			Scope: di.App,
			Build: func(ctn di.Container) (interface{}, error) {
				middle := cfg.Resources.Get("authorize-middleware").(func(next http.Handler) http.Handler)
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
			Name:  "user-storage",
			Scope: di.App,
			Build: func(ctn di.Container) (interface{}, error) {
				database := cfg.Resources.Get("database").(*gorm.DB)
				s := user.NewPGStorage(database)
				s.Drop()
				s.Migrate()
				return s, nil
			},
		},
		{
			Name:  "jwt-service",
			Scope: di.App,
			Build: func(ctn di.Container) (interface{}, error) {
				duration := time.Duration(cfg.JWTExpirationDelta) * time.Second
				return token.NewJWTService(cfg.JWTSigningKey, duration), nil
			},
		},
		{
			Name:  "authenticating-routes",
			Scope: di.App,
			Build: func(ctn di.Container) (interface{}, error) {
				tokenService := cfg.Resources.Get("jwt-service").(authenticating.TokenService)
				store := cfg.Resources.Get("user-storage").(authenticating.Store)
				s := authenticating.NewService(tokenService, store)
				return rest.Authenticating(s), nil
			},
		},
		{
			Name:  "authorize-middleware",
			Scope: di.App,
			Build: func(ctn di.Container) (interface{}, error) {
				tokenService := cfg.Resources.Get("jwt-service").(authorizing.TokenService)
				store := cfg.Resources.Get("user-storage").(authorizing.Store)
				s := authorizing.NewService(tokenService, store)
				return middleware.AuthorizeReq(s), nil
			},
		},
		{
			Name:  "sign_up-routes",
			Scope: di.App,
			Build: func(ctn di.Container) (interface{}, error) {
				store := cfg.Resources.Get("user-storage").(signup.Store)
				s := signup.NewService(store)
				return rest.SignUp(s), nil
			},
		},
		{
			Name:  "repository-storage",
			Scope: di.App,
			Build: func(ctn di.Container) (interface{}, error) {
				database := cfg.Resources.Get("database").(*gorm.DB)
				s := repo.NewPGStorage(database)
				s.Drop()
				s.Migrate()
				return s, nil
			},
		},
		{
			Name:  "creating-routes",
			Scope: di.App,
			Build: func(ctn di.Container) (interface{}, error) {
				store := cfg.Resources.Get("repository-storage").(creating.Store)
				s := creating.NewService(store)
				return rest.Creating(s), nil
			},
		},
		{
			Name:  "listing-services",
			Scope: di.App,
			Build: func(ctn di.Container) (interface{}, error) {
				store := cfg.Resources.Get("repository-storage").(listing.Store)
				return listing.NewService(store), nil
			},
		},
		{
			Name:  "listing-user-repo-routes",
			Scope: di.App,
			Build: func(ctn di.Container) (interface{}, error) {
				s := cfg.Resources.Get("listing-services").(listing.Service)
				return rest.ListingUserRepos(s), nil
			},
		},
		{
			Name:  "listing-public-repos-routes",
			Scope: di.App,
			Build: func(ctn di.Container) (interface{}, error) {
				s := cfg.Resources.Get("listing-services").(listing.Service)
				return rest.ListingPublicRepos(s), nil
			},
		},
	}

	builder.Add(definitions...)
	return builder.Build()
}
