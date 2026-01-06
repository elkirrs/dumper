package runner

import (
	"context"
	"fmt"
)

func RunWithCtx(ctx context.Context, f func() error) error {
	done := make(chan error, 1)
	go func() {
		done <- f()
	}()

	select {
	case <-ctx.Done():
		return fmt.Errorf("operation cancelled")
	case err := <-done:
		return err
	}
}
