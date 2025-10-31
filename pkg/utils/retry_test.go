package utils

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

func TestWithRetry(t *testing.T) {
	t.Run("Success on first try", func(t *testing.T) {
		var calls int32
		err := WithRetry(
			context.Background(),
			3,
			func() error {
				atomic.AddInt32(&calls, 1)
				return nil
			},
			nil,
			nil,
		)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if calls != 1 {
			t.Errorf("expected 1 call, got %d", calls)
		}
	})

	t.Run("Retry until success", func(t *testing.T) {
		orig := BackoffFunc
		BackoffFunc = func(_ int) time.Duration { return 1 * time.Millisecond }
		defer func() { BackoffFunc = orig }()

		var calls int32
		err := WithRetry(
			context.Background(),
			3,
			func() error {
				c := atomic.AddInt32(&calls, 1)
				if c < 3 {
					return errors.New("temporary error")
				}
				return nil
			},
			func(err error) bool { return true },
			nil,
		)
		if err != nil {
			t.Fatalf("expected success, got %v", err)
		}
		if calls != 3 {
			t.Errorf("expected 3 attempts, got %d", calls)
		}
	})

	t.Run("Stops on non-retryable error", func(t *testing.T) {
		var calls int32
		err := WithRetry(
			context.Background(),
			5,
			func() error {
				atomic.AddInt32(&calls, 1)
				return errors.New("fatal")
			},
			func(err error) bool { return false },
			nil,
		)
		if err == nil || err.Error() != "fatal" {
			t.Fatalf("expected fatal error, got %v", err)
		}
		if calls != 1 {
			t.Errorf("expected 1 call, got %d", calls)
		}
	})

	t.Run("Context cancelled", func(t *testing.T) {
		orig := BackoffFunc
		BackoffFunc = func(_ int) time.Duration { return 50 * time.Millisecond }
		defer func() { BackoffFunc = orig }()

		ctx, cancel := context.WithCancel(context.Background())
		var calls int32

		err := WithRetry(
			ctx,
			3,
			func() error {
				atomic.AddInt32(&calls, 1)
				cancel() // отменяем сразу
				return errors.New("retryable error")
			},
			func(err error) bool { return true },
			nil,
		)

		if err == nil || !errors.Is(err, context.Canceled) {
			t.Fatalf("expected context.Canceled, got %v", err)
		}
		if calls != 1 {
			t.Errorf("expected 1 call, got %d", calls)
		}
	})

	t.Run("Calls onRetry correctly", func(t *testing.T) {
		orig := BackoffFunc
		BackoffFunc = func(_ int) time.Duration { return 1 * time.Millisecond }
		defer func() { BackoffFunc = orig }()

		var calls int32
		var retries []int

		err := WithRetry(
			context.Background(),
			3,
			func() error {
				atomic.AddInt32(&calls, 1)
				return errors.New("retryable")
			},
			func(err error) bool { return true },
			func(attempt int, err error) {
				retries = append(retries, attempt)
			},
		)
		if err == nil {
			t.Fatal("expected error after retries")
		}
		if len(retries) != 3 {
			t.Errorf("expected 3 retries, got %v", retries)
		}
	})
}
