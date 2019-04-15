package config

import (
	"time"

	"github.com/go-pg/pg"
	"github.com/pkg/errors"
	"github.com/sarulabs/di"

	"github.com/ifreddyrondon/capture/pkg/adding"
	"github.com/ifreddyrondon/capture/pkg/authenticating"
	"github.com/ifreddyrondon/capture/pkg/authorizing"
	"github.com/ifreddyrondon/capture/pkg/creating"
	"github.com/ifreddyrondon/capture/pkg/getting"
	"github.com/ifreddyrondon/capture/pkg/listing"
	"github.com/ifreddyrondon/capture/pkg/removing"
	"github.com/ifreddyrondon/capture/pkg/signup"
	"github.com/ifreddyrondon/capture/pkg/storage/postgres/capture"
	"github.com/ifreddyrondon/capture/pkg/storage/postgres/repo"
	"github.com/ifreddyrondon/capture/pkg/storage/postgres/user"
	"github.com/ifreddyrondon/capture/pkg/token"
	"github.com/ifreddyrondon/capture/pkg/updating"
)

func getResources(cfg *Config) di.Container {
	builder, _ := di.NewBuilder()
	definitions := []di.Def{
		{
			Name: "database",
			Build: func(ctn di.Container) (i interface{}, e error) {
				opts, err := pg.ParseURL(cfg.Constants.PG)
				if err != nil {
					return nil, errors.Wrap(err, "parsing postgres url into pg options when di")
				}
				db := pg.Connect(opts)
				if db == nil {
					return nil, errors.New("failed to connect to postgres db when di")
				}
				return db, nil
			},
			Close: func(obj interface{}) error {
				return obj.(*pg.DB).Close()
			},
		},
		{
			Name: "user-storage",
			Build: func(ctn di.Container) (interface{}, error) {
				database := cfg.Resources.Get("database").(*pg.DB)
				s := user.NewPGStorage(database)
				if err := s.Drop(); err != nil {
					return nil, errors.Wrap(err, "di dropping schema for user-storage")
				}
				if err := s.CreateSchema(); err != nil {
					return nil, errors.Wrap(err, "di creating schema for user-storage")
				}
				return s, nil
			},
		},
		{
			Name: "sign_up-service",
			Build: func(ctn di.Container) (interface{}, error) {
				store := cfg.Resources.Get("user-storage").(signup.Store)
				return signup.NewService(store), nil
			},
		},
		{
			Name: "jwt-service",
			Build: func(ctn di.Container) (interface{}, error) {
				duration := time.Duration(cfg.JWTExpirationDelta) * time.Second
				return token.NewJWTService(cfg.JWTSigningKey, duration), nil
			},
		},
		{
			Name: "authenticating-service",
			Build: func(ctn di.Container) (interface{}, error) {
				tokenService := cfg.Resources.Get("jwt-service").(authenticating.TokenService)
				store := cfg.Resources.Get("user-storage").(authenticating.Store)
				return authenticating.NewService(tokenService, store), nil
			},
		},
		{
			Name: "authorize-service",
			Build: func(ctn di.Container) (interface{}, error) {
				tokenService := cfg.Resources.Get("jwt-service").(authorizing.TokenService)
				store := cfg.Resources.Get("user-storage").(authorizing.Store)
				return authorizing.NewService(tokenService, store), nil
			},
		},
		{
			Name: "repository-storage",
			Build: func(ctn di.Container) (interface{}, error) {
				database := cfg.Resources.Get("database").(*pg.DB)
				s := repo.NewPGStorage(database)
				if err := s.Drop(); err != nil {
					return nil, errors.Wrap(err, "di dropping schema for repository-storage")
				}
				if err := s.CreateSchema(); err != nil {
					return nil, errors.Wrap(err, "di creating schema for repository-storage")
				}
				return s, nil
			},
		},
		{
			Name: "creating-repo-service",
			Build: func(ctn di.Container) (interface{}, error) {
				store := cfg.Resources.Get("repository-storage").(creating.Store)
				return creating.NewService(store), nil
			},
		},
		{
			Name: "listing-repo-services",
			Build: func(ctn di.Container) (interface{}, error) {
				store := cfg.Resources.Get("repository-storage").(listing.RepoStore)
				return listing.NewRepoService(store), nil
			},
		},
		{
			Name: "getting-repo-service",
			Build: func(ctn di.Container) (interface{}, error) {
				store := cfg.Resources.Get("repository-storage").(getting.RepoStore)
				return getting.NewRepoService(store), nil
			},
		},
		{
			Name: "capture-storage",
			Build: func(ctn di.Container) (interface{}, error) {
				database := cfg.Resources.Get("database").(*pg.DB)
				s := capture.NewPGStorage(database)
				if err := s.Drop(); err != nil {
					return nil, errors.Wrap(err, "di dropping schema for capture-storage")
				}
				if err := s.CreateSchema(); err != nil {
					return nil, errors.Wrap(err, "di creating schema for capture-storage")
				}
				return s, nil
			},
		},
		{
			Name: "adding-capture-service",
			Build: func(ctn di.Container) (interface{}, error) {
				store := cfg.Resources.Get("capture-storage").(adding.CaptureStore)
				return adding.NewCaptureService(store), nil
			},
		},
		{
			Name: "adding-multi-capture-service",
			Build: func(ctn di.Container) (interface{}, error) {
				store := cfg.Resources.Get("capture-storage").(adding.MultiCaptureStore)
				return adding.NewMultiCaptureService(store), nil
			},
		},
		{
			Name: "listing-capture-services",
			Build: func(ctn di.Container) (interface{}, error) {
				store := cfg.Resources.Get("capture-storage").(listing.CaptureStore)
				return listing.NewCaptureService(store), nil
			},
		},
		{
			Name: "getting-capture-service",
			Build: func(ctn di.Container) (interface{}, error) {
				store := cfg.Resources.Get("capture-storage").(getting.CaptureStore)
				return getting.NewCaptureService(store), nil
			},
		},
		{
			Name: "removing-capture-service",
			Build: func(ctn di.Container) (interface{}, error) {
				store := cfg.Resources.Get("capture-storage").(removing.CaptureStore)
				return removing.NewCaptureService(store), nil
			},
		},
		{
			Name: "updating-capture-service",
			Build: func(ctn di.Container) (interface{}, error) {
				store := cfg.Resources.Get("capture-storage").(updating.CaptureStore)
				return updating.NewCaptureService(store), nil
			},
		},
	}

	builder.Add(definitions...)
	return builder.Build()
}
