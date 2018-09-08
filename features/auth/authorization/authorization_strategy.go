package authorization

import (
	"net/http"
)

// Strategy is an Authorization mechanism to validate if the request can access a resource
type Strategy interface {
	// IsAuthorizedREQ validates if a request contains a valid credential.
	IsAuthorizedREQ(*http.Request) (string, error)
	// IsNotAuthorizedErr check if an error is for invalid credentials.
	IsNotAuthorizedErr(error) bool
}
