package rest

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/ifreddyrondon/capture/pkg/adding"
	"github.com/ifreddyrondon/capture/pkg/authenticating"
	"github.com/ifreddyrondon/capture/pkg/authorizing"
	"github.com/ifreddyrondon/capture/pkg/creating"
	"github.com/ifreddyrondon/capture/pkg/getting"
	"github.com/ifreddyrondon/capture/pkg/http/rest/handler"
	"github.com/ifreddyrondon/capture/pkg/http/rest/middleware"
	"github.com/ifreddyrondon/capture/pkg/listing"
	"github.com/ifreddyrondon/capture/pkg/removing"
	"github.com/ifreddyrondon/capture/pkg/signup"
	"github.com/ifreddyrondon/capture/pkg/updating"
	"github.com/sarulabs/di"
)

// Router returns a configured http.Handler with app resources.
func Router(resources di.Container) http.Handler {
	r := chi.NewRouter()

	signUpService := resources.Get("sign_up-service").(signup.Service)
	signUpHandler := handler.SignUp(signUpService)
	authorizeService := resources.Get("authorize-service").(authorizing.Service)
	authorizeMiddleware := middleware.AuthorizeReq(authorizeService)
	authenticatingService := resources.Get("authenticating-service").(authenticating.Service)
	authenticatingHandler := handler.Authenticating(authenticatingService)

	creatingRepoService := resources.Get("creating-repo-service").(creating.Service)
	creatingRepoHandler := handler.Creating(creatingRepoService)
	listingUserReposMiddleware := middleware.FilterUserRepos()
	listingRepoService := resources.Get("listing-repo-services").(listing.RepoService)
	listingUserReposHandler := handler.ListingUserRepos(listingRepoService)
	listingPublicReposMiddleware := middleware.FilterPublicRepos()
	listingPublicReposHandler := handler.ListingPublicRepos(listingRepoService)
	gettingRepoService := resources.Get("getting-repo-service").(getting.RepoService)
	ctxRepoMiddleware := middleware.RepoCtx(gettingRepoService)
	repoOwnerOrPublicMiddleware := middleware.RepoOwnerOrPublic()
	repoOwnerMiddleware := middleware.RepoOwner()
	gettingRepoHandler := handler.GettingRepo()

	addingCaptureService := resources.Get("adding-capture-service").(adding.CaptureService)
	addingCaptureHandler := handler.AddingCapture(addingCaptureService)
	addingMultiCaptureService := resources.Get("adding-multi-capture-service").(adding.MultiCaptureService)
	addingMultiCaptureHandler := handler.AddingMultiCapture(addingMultiCaptureService)
	listingCapturesMiddleware := middleware.FilterCaptures()
	listingCaptureService := resources.Get("listing-capture-services").(listing.CaptureService)
	listingCapturesHandler := handler.ListingRepoCaptures(listingCaptureService)
	gettingCaptureService := resources.Get("getting-capture-service").(getting.CaptureService)
	ctxCaptureMiddleware := middleware.CaptureCtx(gettingCaptureService)
	gettingCaptureHandler := handler.GettingCapture()
	removingCaptureService := resources.Get("removing-capture-service").(removing.CaptureService)
	removingCaptureHandler := handler.RemovingCapture(removingCaptureService)
	updatingCaptureService := resources.Get("updating-capture-service").(updating.CaptureService)
	updatingCaptureHandler := handler.UpdatingCapture(updatingCaptureService)

	r.Post("/sign/", signUpHandler)
	r.Route("/auth/", func(r chi.Router) {
		r.Post("/token-auth", authenticatingHandler)
	})
	r.Route("/user/", func(r chi.Router) {
		r.Use(authorizeMiddleware)
		r.Route("/repos/", func(r chi.Router) {
			r.Post("/", creatingRepoHandler)
			r.With(listingUserReposMiddleware).Get("/", listingUserReposHandler)

		})
	})
	r.Route("/repositories/", func(r chi.Router) {
		r.Use(authorizeMiddleware)
		r.With(listingPublicReposMiddleware).
			Get("/", listingPublicReposHandler)
		r.Route("/{id}", func(r chi.Router) {
			r.Use(ctxRepoMiddleware)
			r.With(repoOwnerOrPublicMiddleware).Get("/", gettingRepoHandler)
			r.Route("/captures/", func(r chi.Router) {
				r.With(repoOwnerMiddleware).Post("/", addingCaptureHandler)
				r.With(repoOwnerMiddleware).Post("/multi", addingMultiCaptureHandler)
				r.With(repoOwnerOrPublicMiddleware).With(listingCapturesMiddleware).Get("/", listingCapturesHandler)
				r.Route("/{captureId}", func(r chi.Router) {
					r.Use(ctxCaptureMiddleware)
					r.With(repoOwnerOrPublicMiddleware).Get("/", gettingCaptureHandler)
					r.With(repoOwnerMiddleware).Delete("/", removingCaptureHandler)
					r.With(repoOwnerMiddleware).Put("/", updatingCaptureHandler)
				})
			})
		})
	})

	return r
}
