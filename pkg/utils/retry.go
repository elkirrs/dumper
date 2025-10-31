package utils

import (
	"context"
	"dumper/pkg/logging"
	"fmt"
	"time"
)

var BackoffFunc = ExponentialBackoff

func WithRetry(
	ctx context.Context,
	maxRetries int,
	fn func() error,
	shouldRetry func(err error) bool,
	onRetry func(attempt int, err error),
) error {
	var err error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		err = fn()
		if err == nil {
			return nil
		}

		if shouldRetry != nil && shouldRetry(err) {
			if onRetry != nil {
				onRetry(attempt, err)
			}

			if attempt == maxRetries {
				logging.L(ctx).Error("Failed retrying connection",
					logging.IntAttr("attempts", maxRetries),
					logging.ErrAttr(err),
				)
				return fmt.Errorf("failed after %d attempts: %w", maxRetries, err)
			}

			delay := ExponentialBackoff(attempt)

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
