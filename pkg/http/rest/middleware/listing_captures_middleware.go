package middleware

import (
	"net/http"

	"github.com/ifreddyrondon/bastion/middleware"
)

const capturesMaxAllowedLimit = 100

func FilterCaptures() func(next http.Handler) http.Handler {
	return middleware.Listing(
		middleware.MaxAllowedLimit(capturesMaxAllowedLimit),
		middleware.Sort(updatedDESC, updatedASC, createdDESC, createdASC),
	)
}
