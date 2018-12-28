package adding

import (
	"time"
)

func SetPostTimestampInstance(targetTimestamp *Timestamp, t time.Time) {
	targetTimestamp.postTimestamp = &t
}
