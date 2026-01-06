package runner

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestRunWithCtx(t *testing.T) {
	t.Run("Function completes successfully", func(t *testing.T) {
		ctx := context.Background()
		err := RunWithCtx(ctx, func() error {
			time.Sleep(10 * time.Millisecond)
			return nil
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("Function returns error", func(t *testing.T) {
		ctx := context.Background()
		expected := errors.New("something went wrong")
		err := RunWithCtx(ctx, func() error {
			return expected
		})
		if !errors.Is(err, expected) {
			t.Fatalf("expected %v, got %v", expected, err)
		}
	})

	t.Run("Context canceled before function finishes", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
		defer cancel()

		err := RunWithCtx(ctx, func() error {
			time.Sleep(200 * time.Millisecond)
			return nil
		})
		if err == nil || err.Error() != "operation cancelled" {
			t.Fatalf("expected operation cancelled, got %v", err)
		}
	})

	t.Run("Context already canceled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := RunWithCtx(ctx, func() error {
			time.Sleep(100 * time.Millisecond)
			return nil
		})
		if err == nil || err.Error() != "operation cancelled" {
			t.Fatalf("expected operation cancelled, got %v", err)
		}
	})

	t.Run("Function finishes before context cancel", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		err := RunWithCtx(ctx, func() error {
			time.Sleep(20 * time.Millisecond)
			return nil
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})
}
