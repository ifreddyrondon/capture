package authentication

import (
	"github.com/ifreddyrondon/capture/app/user"
)

// Strategy is an Authentication mechanisms to validate users credentials
type Strategy interface {
	// Validate user credentials from bytes.
	Validate([]byte) (*user.User, error)
	// IsErrCredentials check if an error is for invalid credentials.
	IsErrCredentials(error) bool
	// IsErrDecoding check if an error is for invalid decoding credentials.
	IsErrDecoding(error) bool
}
