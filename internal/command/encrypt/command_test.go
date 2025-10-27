package encrypt

import (
	"dumper/internal/domain/encrypt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncrypt_Generate_AES(t *testing.T) {
	tests := []struct {
		name        string
		opts        *encrypt.Options
		wantCMD     string
		wantName    string
		expectError bool
	}{
		{
			name: "AES encrypt",
			opts: &encrypt.Options{
				FilePath: "testfile.sql",
				Password: "123456",
				Type:     "aes",
				Crypt:    "encrypt",
			},
			wantCMD:  "openssl enc -aes-256-cbc -salt -pbkdf2 -iter 100000 -in testfile.sql -out testfile.sql.enc -k 123456",
			wantName: "testfile.sql.enc",
		},
		{
			name: "AES decrypt",
			opts: &encrypt.Options{
				FilePath: "testfile.sql.enc",
				Password: "123456",
				Type:     "aes",
				Crypt:    "decrypt",
			},
			wantCMD:  "openssl enc -d -aes-256-cbc -pbkdf2 -iter 100000 -in testfile.sql.enc -out testfile.sql -k 123456",
			wantName: "testfile.sql",
		},
		{
			name: "Unknown type",
			opts: &encrypt.Options{
				FilePath: "testfile.sql",
				Type:     "unknown",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := NewApp(tt.opts)
			got, err := app.Generate()

			if tt.expectError {
				assert.Error(t, err, "expected an error for unknown encrypt type")
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.wantCMD, got.CMD, "generated command should match expected")
			assert.Equal(t, tt.wantName, got.Name, "generated filename should match expected")
		})
	}
}
