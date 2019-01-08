package middleware

import (
	"net/http"

	"github.com/ifreddyrondon/bastion/middleware"
	"github.com/ifreddyrondon/bastion/middleware/listing/filtering"
)

const reposMaxAllowedLimit = 50

func FilterOwnRepos() func(next http.Handler) http.Handler {
	publicVisibility := filtering.NewValue("public", "public repos")
	privateVisibility := filtering.NewValue("private", "private repos")
	visibilityFilter := filtering.NewText("visibility", "filters the repos by their visibility", publicVisibility, privateVisibility)

	return middleware.Listing(
		middleware.MaxAllowedLimit(reposMaxAllowedLimit),
		middleware.Sort(updatedDESC, updatedASC, createdDESC, createdASC),
		middleware.Filter(visibilityFilter),
	)
}

func FilterPublicRepos() func(next http.Handler) http.Handler {
	return middleware.Listing(
		middleware.MaxAllowedLimit(reposMaxAllowedLimit),
		middleware.Sort(updatedDESC, updatedASC, createdDESC, createdASC),
	)
}
