package aes

import (
	"dumper/internal/domain/encrypt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAESGenerator_Generate(t *testing.T) {
	tests := []struct {
		name     string
		opts     *encrypt.Options
		wantCMD  string
		wantName string
	}{
		{
			name: "encrypt file",
			opts: &encrypt.Options{
				FilePath: "testfile.sql",
				Password: "123456",
				Crypt:    "encrypt",
			},
			wantCMD:  "openssl enc -aes-256-cbc -salt -pbkdf2 -iter 100000 -in testfile.sql -out testfile.sql.enc -k 123456",
			wantName: "testfile.sql.enc",
		},
		{
			name: "decrypt file",
			opts: &encrypt.Options{
				FilePath: "testfile.sql.enc",
				Password: "123456",
				Crypt:    "decrypt",
			},
			wantCMD:  "openssl enc -d -aes-256-cbc -pbkdf2 -iter 100000 -in testfile.sql.enc -out testfile.sql -k 123456",
			wantName: "testfile.sql",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen := AESGenerator{}
			got := gen.Generate(tt.opts)

			assert.Equal(t, tt.wantCMD, got.CMD, "generated command should match expected")
			assert.Equal(t, tt.wantName, got.Name, "generated output filename should match expected")
		})
	}
}
