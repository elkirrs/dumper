package progress_test

import (
	"bytes"
	"dumper/pkg/utils/progress"
	"io"
	"os"
	"regexp"
	"testing"
)

func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	f()
	_ = w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	return buf.String()
}

func normalizeOutput(s string) string {
	s = regexp.MustCompile(`\$begin:math:display\$|\$end:math:display\$`).ReplaceAllString(s, "")
	return s
}

func TestProgress(t *testing.T) {
	tests := []struct {
		name     string
		done     int64
		total    int64
		expected *regexp.Regexp
	}{
		{
			name:     "Zero total",
			done:     512,
			total:    512,
			expected: regexp.MustCompile(`Uploaded: 100\.0% \[512/512 bytes\]`),
		},
		{
			name:     "Half progress",
			done:     50,
			total:    100,
			expected: regexp.MustCompile(`Uploading\.\.\. 50\.0% \[50/100 bytes\]`),
		},
		{
			name:     "Full progress",
			done:     200,
			total:    200,
			expected: regexp.MustCompile(`Uploaded: 100\.0% \[200/200 bytes\]`),
		},
		{
			name:     "Non integer percentage",
			done:     1,
			total:    3,
			expected: regexp.MustCompile(`Uploading\.\.\. 33\.3% \[1/3 bytes\]`),
		},
		{
			name:     "Over total (more than 100%)",
			done:     120,
			total:    100,
			expected: regexp.MustCompile(`Uploading\.\.\. 120\.0% \[120/100 bytes\]`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := captureOutput(func() {
				progress.Progress(tt.done, tt.total)
			})
			out = normalizeOutput(out)

			if !tt.expected.MatchString(out) {
				t.Errorf("\nexpected match:\n%v\nbut got:\n%q", tt.expected, out)
			}
		})
	}
}
