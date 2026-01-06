package attempt

import (
	"math/rand"
	"time"
)

func ExponentialBackoff(attempt int) time.Duration {
	base := time.Duration(1<<uint(attempt-1)) * time.Second
	jitterRange := int64(float64(base) * 0.6)
	jitter := time.Duration(rand.Int63n(jitterRange) - jitterRange/2)
	delay := base + jitter
	if delay < 0 {
		delay = 0
	}
	return delay
}
