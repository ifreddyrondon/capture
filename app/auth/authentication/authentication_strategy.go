package authentication

import (
	"net/http"

	"github.com/ifreddyrondon/capture/app/user"
)

// Strategy is an Authentication mechanisms to validate users credentials
type Strategy interface {
	// Validate user credentials from bytes.
	Validate(*http.Request) (*user.User, error)
	// IsErrCredentials check if an error is for invalid credentials.
	IsErrCredentials(error) bool
	// IsErrDecoding check if an error is for invalid decoding credentials.
	IsErrDecoding(error) bool
}
