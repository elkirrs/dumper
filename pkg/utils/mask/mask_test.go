package mask_test

import (
	"dumper/pkg/utils/mask"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMask(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"mysql://user:secret@localhost/db", "mysql://********:********@localhost/db"},
		{"user:pass@host", "********:********@host"},
		{"--password=secret", "--password=********"},
		{"PWD='secret'", "PWD='********'"},
		{"--user secret", "--user ********"},
		{"PGPASSWORD secret", "PGPASSWORD ********"},
		{"-psecret", "-p********"},
		{"-uuser", "-u********"},
		{"MYSQL_PWD=mysecret", "MYSQL_PWD=********"},
		{"AWS_SECRET_ACCESS_KEY=abc123", "AWS_SECRET_ACCESS_KEY=********"},
		{"user:pass@tcp(localhost:3306)", "********:********@tcp(localhost:3306)"},
		{`user:pass@host \'test\'`, `********:********@host 'test'`},
		{`user:pass@host \"test\"`, `********:********@host "test"`},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			out := mask.Mask(tt.input)
			assert.Equal(t, tt.expected, out)
		})
	}
}
