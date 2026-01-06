package suffix

import "testing"

func TestRemoveSuffix(t *testing.T) {
	tests := []struct {
		name     string
		str      string
		suffix   string
		expected string
	}{
		{
			name:     "remove valid suffix",
			str:      "filename.sql.gz",
			suffix:   ".gz",
			expected: "filename.sql",
		},
		{
			name:     "no suffix match",
			str:      "filename.sql",
			suffix:   ".gz",
			expected: "filename.sql",
		},
		{
			name:     "empty suffix",
			str:      "filename.sql",
			suffix:   "",
			expected: "filename.sql",
		},
		{
			name:     "suffix equals string",
			str:      "test",
			suffix:   "test",
			expected: "",
		},
		{
			name:     "partial match at end but not suffix",
			str:      "testing",
			suffix:   "sting",
			expected: "te",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RemoveSuffix(tt.str, tt.suffix)
			if got != tt.expected {
				t.Errorf("RemoveSuffix(%q, %q) = %q; want %q", tt.str, tt.suffix, got, tt.expected)
			}
		})
	}
}
