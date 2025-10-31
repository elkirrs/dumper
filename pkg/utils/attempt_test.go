package utils_test

import (
	"dumper/pkg/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestExponentialBackoff_Positive(t *testing.T) {
	for attempt := 1; attempt <= 10; attempt++ {
		delay := utils.ExponentialBackoff(attempt)
		assert.GreaterOrEqual(t, delay, time.Duration(0), "Delay should not be negative")
	}
}

func TestExponentialBackoff_ExponentialGrowth(t *testing.T) {
	attempt1 := utils.ExponentialBackoff(1)
	attempt2 := utils.ExponentialBackoff(2)
	attempt3 := utils.ExponentialBackoff(3)

	assert.LessOrEqual(t, attempt1*2, attempt2+time.Second*2, "Attempt 2 should be roughly double attempt 1")
	assert.LessOrEqual(t, attempt2*2, attempt3+time.Second*4, "Attempt 3 should be roughly double attempt 2")
}

func TestExponentialBackoff_JitterRange(t *testing.T) {
	for attempt := 1; attempt <= 5; attempt++ {
		base := time.Duration(1<<uint(attempt-1)) * time.Second
		jitterRange := int64(float64(base) * 0.6)

		for i := 0; i < 100; i++ {
			delay := utils.ExponentialBackoff(attempt)
			min := base - time.Duration(jitterRange/2)
			max := base + time.Duration(jitterRange/2)
			if min < 0 {
				min = 0
			}
			assert.GreaterOrEqual(t, delay, min, "Delay is below expected range")
			assert.LessOrEqual(t, delay, max, "Delay is above expected range")
		}
	}
}
