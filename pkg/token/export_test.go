package token

import (
	"github.com/ifreddyrondon/capture/pkg"
)

// SetClockInstance is a helper function only exported for test.
// It's intended to be used for stub the time.Now() function.
func SetClockInstance(targetClaims *JWTClaims, clock *pkg.Clock) {
	targetClaims.clock = clock
}
