package app

import (
	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/bastion/render/json"
	"github.com/ifreddyrondon/gocapture/auth"
	"github.com/ifreddyrondon/gocapture/auth/strategy/basic"
	"github.com/ifreddyrondon/gocapture/branch"
	"github.com/ifreddyrondon/gocapture/capture"
	"github.com/ifreddyrondon/gocapture/jwt"
	"github.com/ifreddyrondon/gocapture/user"
	"github.com/jinzhu/gorm"
)

// ContextKey represent a key in a context.ContextValue
type ContextKey string

func (c ContextKey) String() string {
	return string(c)
}

// New returns a bastion ready instance with all the app config
func New(db *gorm.DB) *bastion.Bastion {
	app := bastion.New(bastion.Options{})

	userRepo := user.NewPGRepository(db)
	userRepo.Drop()
	userRepo.Migrate()
	userController := user.NewController(userRepo, json.NewRender)
	app.APIRouter.Mount("/users/", userController.Router())

	strategy := basic.Strategy{
		Render:        json.NewRender,
		UserKey:       ContextKey("user"),
		GetterService: userRepo,
	}

	jwtService := jwt.NewService([]byte("test"), jwt.DefaultJWTExpirationDelta, json.NewRender)

	authController := auth.Controller{
		Strategy: strategy,
		Render:   json.NewRender,
		UserKey:  ContextKey("user"),
		JWT:      jwtService,
	}
	app.APIRouter.Mount("/auth/", authController.Router())

	captureService := capture.PGService{DB: db}
	captureService.Drop()
	captureService.Migrate()
	captureController := capture.Controller{
		Service: &captureService,
		Render:  json.NewRender,
		CtxKey:  ContextKey("capture"),
	}
	app.APIRouter.Mount("/captures/", captureController.Router())

	branchController := branch.Controller{
		Render: json.NewRender,
	}
	app.APIRouter.Mount("/branches/", branchController.Router())
	return app
}
