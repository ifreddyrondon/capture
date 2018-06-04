package auth

import "net/http"

// Strategy is an Authentication mechanisms to validate users credentials
type Strategy interface {
	// Authenticate validate if an user is authorized to continue or 401.
	Authenticate(next http.Handler) http.Handler
}

// Authorization validates if a request is authorized
type Authorization interface {
	IsAuthorized(next http.Handler) http.Handler
}