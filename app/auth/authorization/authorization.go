package authorization

import "net/http"

// Authorization validates if a request is authorized
type Authorization interface {
	IsAuthorized(next http.Handler) http.Handler
}
