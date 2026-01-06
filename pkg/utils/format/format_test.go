package format_test

import (
	"dumper/pkg/utils/format"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		input    int64
		expected string
	}{
		{512, "512 B"},                // < 1 KB
		{1024, "1.00 KB"},             // = 1 KB
		{1536, "1.50 KB"},             // 1.5 KB
		{1048576, "1.00 MB"},          // 1 MB
		{1572864, "1.50 MB"},          // 1.5 MB
		{1073741824, "1.00 GB"},       // 1 GB
		{1610612736, "1.50 GB"},       // 1.5 GB
		{1099511627776, "1.00 TB"},    // 1 TB
		{1125899906842624, "1.00 PB"}, // 1 PB
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := format.FormatBytes(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
