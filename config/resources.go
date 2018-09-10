package config

import (
	"github.com/ifreddyrondon/capture/features/capture"
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
				database := cfg.Container.Get("database").(*gorm.DB)
				store := user.NewPGStore(database)
				store.Drop()
				store.Migrate()
				return user.NewService(store), nil
			},
		},
		{
			Name:  "repo-service",
			Scope: di.App,
			Build: func(ctn di.Container) (interface{}, error) {
				database := cfg.Container.Get("database").(*gorm.DB)
				store := repository.NewPGStore(database)
				store.Drop()
				store.Migrate()
				return repository.Service(store), nil
			},
		},
		{
			Name:  "capture-service",
			Scope: di.App,
			Build: func(ctn di.Container) (interface{}, error) {
				database := cfg.Container.Get("database").(*gorm.DB)
				store := capture.NewPGStore(database)
				store.Drop()
				store.Migrate()
				return capture.Service(store), nil
			},
		},
	}

	builder.Add(definitions...)
	return builder.Build()
}
