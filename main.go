package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ifreddyrondon/bastion/middleware"
	"github.com/ifreddyrondon/bastion/middleware/listing/filtering"
	"github.com/ifreddyrondon/bastion/middleware/listing/sorting"

	"github.com/go-chi/chi"
	"github.com/ifreddyrondon/bastion"
	"github.com/ifreddyrondon/capture/config"
	"github.com/sarulabs/di"
)

func router(resources di.Container) http.Handler {
	r := chi.NewRouter()

	authorize := resources.Get("authorize-middleware").(func(next http.Handler) http.Handler)

	signup := resources.Get("sign_up-routes").(http.HandlerFunc)
	authenticating := resources.Get("authenticating-routes").(http.HandlerFunc)
	creating := resources.Get("creating-routes").(http.HandlerFunc)

	updatedDESC := sorting.NewSort("updated_at_desc", "updated_at DESC", "Updated date descending")
	updatedASC := sorting.NewSort("updated_at_asc", "updated_at ASC", "Updated date ascendant")
	createdDESC := sorting.NewSort("created_at_desc", "created_at DESC", "Created date descending")
	createdASC := sorting.NewSort("created_at_asc", "created_at ASC", "Created date ascendant")

	publicVisibility := filtering.NewValue("public", "public repos")
	privateVisibility := filtering.NewValue("private", "private repos")
	visibilityFilter := filtering.NewText("visibility", "filters the repos by their visibility", publicVisibility, privateVisibility)

	listingMiddleware := middleware.Listing(
		middleware.MaxAllowedLimit(50),
		middleware.Sort(updatedDESC, updatedASC, createdDESC, createdASC),
		middleware.Filter(visibilityFilter),
	)
	listingRepos := resources.Get("listing-repo-routes").(http.HandlerFunc)

	repositoryRoutes := resources.Get("repositories-routes").(http.Handler)
	captureRoutes := resources.Get("capture-routes").(http.Handler)
	branchRoutes := resources.Get("branch-routes").(http.Handler)
	multipostRoutes := resources.Get("multipost-routes").(http.Handler)

	r.Post("/sign/", signup)
	r.Route("/auth/", func(r chi.Router) {
		r.Post("/token-auth", authenticating)
	})
	r.Route("/user/", func(r chi.Router) {
		r.Route("/repos/", func(r chi.Router) {
			r.Use(authorize)
			r.Post("/", creating)
			r.Route("/", func(r chi.Router) {
				r.Use(listingMiddleware)
				r.Get("/", listingRepos)
			})

		})
	})
	r.Mount("/repositories/", repositoryRoutes)
	r.Mount("/captures/", captureRoutes)
	r.Mount("/branches/", branchRoutes)
	r.Mount("/multipost/", multipostRoutes)

	return r
}

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Panicln("Configuration error", err)
	}

	app := bastion.New(bastion.Addr(cfg.ADDR))
	app.APIRouter.Mount("/", router(cfg.Resources))
	app.RegisterOnShutdown(cfg.OnShutdown)
	if err := app.Serve(); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
}
