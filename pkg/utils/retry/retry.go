package retry

import (
	"context"
	"dumper/pkg/logging"
	"dumper/pkg/utils/attempt"
	"fmt"
	"time"
)

var BackoffFunc = attempt.ExponentialBackoff

func WithRetry(
	ctx context.Context,
	maxRetries int,
	fn func() error,
	shouldRetry func(err error) bool,
	onRetry func(attempt int, err error),
) error {
	var err error
	for attempts := 1; attempts <= maxRetries; attempts++ {
		err = fn()
		if err == nil {
			return nil
		}

		if shouldRetry != nil && shouldRetry(err) {
			if onRetry != nil {
				onRetry(attempts, err)
			}

			if attempts == maxRetries {
				logging.L(ctx).Error("Failed retrying connection",
					logging.IntAttr("attempts", maxRetries),
					logging.ErrAttr(err),
				)
				return fmt.Errorf("failed after %d attempts: %w", maxRetries, err)
			}

			delay := attempt.ExponentialBackoff(attempts)

			logging.L(ctx).Error("Connection error, retrying after",
				logging.StringAttr("time", delay.String()),
				logging.ErrAttr(err),
			)
			fmt.Printf("Connection error, retrying after %.2fs\n", delay.Seconds())

			select {
			case <-ctx.Done():
				return fmt.Errorf("operation cancelled: %w", ctx.Err())
			case <-time.After(delay):
				continue
			}
		}

		return err
	}
	return err
}
