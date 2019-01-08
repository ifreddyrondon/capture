package config

import (
	"time"

	"github.com/go-pg/pg"
	"github.com/ifreddyrondon/capture/pkg/adding"
	"github.com/ifreddyrondon/capture/pkg/getting"
	"github.com/ifreddyrondon/capture/pkg/listing"
	"github.com/ifreddyrondon/capture/pkg/removing"
	"github.com/ifreddyrondon/capture/pkg/storage/postgres/capture"
	"github.com/ifreddyrondon/capture/pkg/storage/postgres/repo"
	"github.com/ifreddyrondon/capture/pkg/updating"
	"github.com/pkg/errors"

	"github.com/ifreddyrondon/capture/pkg/creating"

	"github.com/ifreddyrondon/capture/pkg/authenticating"
	"github.com/ifreddyrondon/capture/pkg/authorizing"
	"github.com/ifreddyrondon/capture/pkg/http/rest"
	"github.com/ifreddyrondon/capture/pkg/http/rest/middleware"
	"github.com/ifreddyrondon/capture/pkg/signup"
	"github.com/ifreddyrondon/capture/pkg/storage/postgres/user"
	"github.com/ifreddyrondon/capture/pkg/token"
	"github.com/sarulabs/di"
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
			Name: "jwt-service",
			Build: func(ctn di.Container) (interface{}, error) {
				duration := time.Duration(cfg.JWTExpirationDelta) * time.Second
				return token.NewJWTService(cfg.JWTSigningKey, duration), nil
			},
		},
		{
			Name: "authenticating-routes",
			Build: func(ctn di.Container) (interface{}, error) {
				tokenService := cfg.Resources.Get("jwt-service").(authenticating.TokenService)
				store := cfg.Resources.Get("user-storage").(authenticating.Store)
				s := authenticating.NewService(tokenService, store)
				return rest.Authenticating(s), nil
			},
		},
		{
			Name: "authorize-middleware",
			Build: func(ctn di.Container) (interface{}, error) {
				tokenService := cfg.Resources.Get("jwt-service").(authorizing.TokenService)
				store := cfg.Resources.Get("user-storage").(authorizing.Store)
				s := authorizing.NewService(tokenService, store)
				return middleware.AuthorizeReq(s), nil
			},
		},
		{
			Name: "sign_up-routes",
			Build: func(ctn di.Container) (interface{}, error) {
				store := cfg.Resources.Get("user-storage").(signup.Store)
				s := signup.NewService(store)
				return rest.SignUp(s), nil
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
			Name: "ctx-repo-middleware",
			Build: func(ctn di.Container) (interface{}, error) {
				store := cfg.Resources.Get("repository-storage").(getting.RepoStore)
				s := getting.NewRepoService(store)
				return middleware.RepoCtx(s), nil
			},
		},
		{
			Name: "repo-owner-or-public-middleware",
			Build: func(ctn di.Container) (interface{}, error) {
				return middleware.RepoOwnerOrPublic(), nil
			},
		},
		{
			Name: "repo-owner-middleware",
			Build: func(ctn di.Container) (interface{}, error) {
				return middleware.RepoOwner(), nil
			},
		},
		{
			Name: "creating-repo-routes",
			Build: func(ctn di.Container) (interface{}, error) {
				store := cfg.Resources.Get("repository-storage").(creating.Store)
				s := creating.NewService(store)
				return rest.Creating(s), nil
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
			Name: "listing-user-repos-middleware",
			Build: func(ctn di.Container) (interface{}, error) {
				return middleware.FilterOwnRepos(), nil
			},
		},
		{
			Name: "listing-user-repos-routes",
			Build: func(ctn di.Container) (interface{}, error) {
				s := cfg.Resources.Get("listing-repo-services").(listing.RepoService)
				return rest.ListingUserRepos(s), nil
			},
		},
		{
			Name: "listing-public-repos-middleware",
			Build: func(ctn di.Container) (interface{}, error) {
				return middleware.FilterPublicRepos(), nil
			},
		},
		{
			Name: "listing-public-repos-routes",
			Build: func(ctn di.Container) (interface{}, error) {
				s := cfg.Resources.Get("listing-repo-services").(listing.RepoService)
				return rest.ListingPublicRepos(s), nil
			},
		},
		{
			Name: "getting-repo-routes",
			Build: func(ctn di.Container) (interface{}, error) {
				return rest.GettingRepo(), nil
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
			Name: "adding-capture-routes",
			Build: func(ctn di.Container) (interface{}, error) {
				store := cfg.Resources.Get("capture-storage").(adding.CaptureStore)
				s := adding.NewCaptureService(store)
				return rest.AddingCapture(s), nil
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
			Name: "listing-captures-middleware",
			Build: func(ctn di.Container) (interface{}, error) {
				return middleware.FilterCaptures(), nil
			},
		},
		{
			Name: "listing-captures-routes",
			Build: func(ctn di.Container) (interface{}, error) {
				s := cfg.Resources.Get("listing-capture-services").(listing.CaptureService)
				return rest.ListingRepoCaptures(s), nil
			},
		},
		{
			Name: "ctx-capture-middleware",
			Build: func(ctn di.Container) (interface{}, error) {
				store := cfg.Resources.Get("capture-storage").(getting.CaptureStore)
				s := getting.NewCaptureService(store)
				return middleware.CaptureCtx(s), nil
			},
		},
		{
			Name: "getting-capture-routes",
			Build: func(ctn di.Container) (interface{}, error) {
				return rest.GettingCapture(), nil
			},
		},
		{
			Name: "removing-capture-routes",
			Build: func(ctn di.Container) (interface{}, error) {
				store := cfg.Resources.Get("capture-storage").(removing.CaptureStore)
				s := removing.NewCaptureService(store)
				return rest.RemovingCapture(s), nil
			},
		},
		{
			Name: "updating-capture-routes",
			Build: func(ctn di.Container) (interface{}, error) {
				store := cfg.Resources.Get("capture-storage").(updating.CaptureStore)
				s := updating.NewCaptureService(store)
				return rest.UpdatingCapture(s), nil
			},
		},
	}

	builder.Add(definitions...)
	return builder.Build()
}
