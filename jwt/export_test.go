package jwt

import (
	"github.com/ifreddyrondon/gocapture/timestamp"
)

// SetClockInstance is a helper function only exported for test.
// It's intended to be used for stub the time.Now() function.
func SetClockInstance(targetClaims *Claims, clock *timestamp.Clock) {
	targetClaims.clock = clock
}
