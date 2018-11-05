package config

import (
	"github.com/ifreddyrondon/capture/features/auth"
	"github.com/ifreddyrondon/capture/features/auth/authentication"
	"github.com/ifreddyrondon/capture/features/auth/authentication/strategy/basic"
	"github.com/ifreddyrondon/capture/features/auth/authorization"
	"github.com/ifreddyrondon/capture/features/auth/jwt"
	"github.com/ifreddyrondon/capture/features/branch"
	"github.com/ifreddyrondon/capture/features/capture"
	"github.com/ifreddyrondon/capture/features/multipost"
	"github.com/ifreddyrondon/capture/features/repository"
	"github.com/ifreddyrondon/capture/features/user"
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
				service := capture.NewService(store)
				return capture.Routes(service), nil
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
			Name:  "repo-routes",
			Scope: di.App,
			Build: func(ctn di.Container) (interface{}, error) {
				userService := cfg.Resources.Get("user-service").(user.Service)
				loggedUser := user.LoggedUser(userService)
				jwtService := cfg.Resources.Get("jwt-service").(*jwt.Service)
				isAuth := authorization.IsAuthorizedREQ(jwtService)
				database := cfg.Resources.Get("database").(*gorm.DB)
				store := repository.NewPGStore(database)
				store.Drop()
				store.Migrate()
				return repository.Routes(store, isAuth, loggedUser), nil
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
