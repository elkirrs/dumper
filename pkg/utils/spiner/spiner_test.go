package spiner

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"time"
)

func TestSpinner(t *testing.T) {
	origStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	stop := make(chan struct{})
	done := make(chan struct{})

	go func() {
		defer close(done)
		Spinner(stop)
	}()

	time.Sleep(500 * time.Millisecond)
	close(stop)

	select {
	case <-done:
	case <-time.After(1 * time.Second):
		t.Fatal("spinner did not stop within timeout")
	}

	_ = w.Close()
	os.Stdout = origStdout

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	output := buf.String()

	if !strings.Contains(output, "Creating dump") {
		t.Errorf("expected output to contain 'Creating dump', got: %q", output)
	}

	if !strings.ContainsAny(output, "-\\|/") {
		t.Errorf("expected spinner characters (-\\|/), got: %q", output)
	}

	if !strings.Contains(output, "\r") {
		t.Errorf("expected carriage return in output, got: %q", output)
	}

	fmt.Println("Spinner test output sample:", output)
}
