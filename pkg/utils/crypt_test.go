package utils_test

import (
	"dumper/internal/domain/app"
	"dumper/pkg/utils"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncryptDecryptAES(t *testing.T) {
	key := []byte("12345678901234567890123456789012")
	plain := []byte("hello world")

	encrypted := utils.EncryptAES(key, plain)
	assert.NotEqual(t, plain, encrypted)

	decrypted, err := utils.DecryptAES(key, encrypted)
	assert.NoError(t, err)
	assert.Equal(t, plain, decrypted)
}

func TestDecryptAES_InvalidCiphertext(t *testing.T) {
	key := []byte("12345678901234567890123456789012")

	_, err := utils.DecryptAES(key, []byte("short"))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ciphertext too short")
}

func TestDeriveKey_Consistency(t *testing.T) {
	password := "pass"
	appSecret := []byte("secret")
	salt := []byte("salt")

	key1 := utils.DeriveKey(password, appSecret, salt)
	key2 := utils.DeriveKey(password, appSecret, salt)

	assert.Equal(t, key1, key2)
}

func TestDeriveAppKey_Consistency(t *testing.T) {
	appSecret := []byte("secret")
	deviceID := []byte("device")

	key := utils.DeriveAppKey(appSecret, deviceID)
	assert.Len(t, key, 32)
}

func TestComputeFinalKey_Consistency(t *testing.T) {
	passwordKey := []byte("passwordkey1234567890123456")
	deviceKey := []byte("devicekey1234567890123456789")

	key1 := utils.ComputeFinalKey(passwordKey, deviceKey)
	key2 := utils.ComputeFinalKey(passwordKey, deviceKey)

	assert.Equal(t, key1, key2)
}

func TestMagicHeaderAndGetScope(t *testing.T) {
	tests := []struct {
		scope   string
		wantHdr string
	}{
		{"app", "ENCADCA"},
		{"device", "ENCDDCA"},
		{"both", "ENCFDCA"},
	}

	for _, tt := range tests {
		hdr := utils.MagicHeader(tt.scope)
		assert.Equal(t, tt.wantHdr, string(hdr))
		assert.Equal(t, tt.scope, utils.GetScope(string(hdr)))
	}
}

func TestIsEncrypted(t *testing.T) {
	h := utils.MagicHeader("app")
	data := append(h, []byte("payload")...)

	assert.True(t, utils.IsEncrypted(data))
	assert.False(t, utils.IsEncrypted([]byte("not_encrypted_data")))
}

func TestLooksEncrypted(t *testing.T) {
	d := make([]byte, 32)
	for i := 0; i < 16; i++ {
		d[i] = 200
	}
	assert.True(t, utils.LooksEncrypted(d))

	ascii := []byte("this is ascii data.......................")
	assert.False(t, utils.LooksEncrypted(ascii))
}

func TestSecretCrypt(t *testing.T) {
	tests := []struct {
		name      string
		scope     string
		expectApp bool
		expectDev bool
	}{
		{"app only", "app", true, false},
		{"device only", "device", false, true},
		{"both", "both", true, true},
		{"default", "", true, true}, // fallback
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flags := &app.Flags{
				Scope:     tt.scope,
				AppSecret: "APPSECRET",
			}

			s := utils.SecretCrypt(flags)

			if tt.expectApp {
				assert.Equal(t, []byte("APPSECRET"), s.SecretKey)
			} else {
				assert.Nil(t, s.SecretKey)
			}

			if tt.expectDev {
				assert.NotNil(t, s.DeviceKey)
				assert.NotEmpty(t, s.DeviceKey)
			} else {
				assert.Nil(t, s.DeviceKey)
			}
		})
	}
}

func TestIsEncryptedFileAndReadEncryptedFile_AllScopes(t *testing.T) {
	scopes := []string{"app", "device", "both"}

	for _, scope := range scopes {
		t.Run("scope="+scope, func(t *testing.T) {
			tmpDir := t.TempDir()
			path := filepath.Join(tmpDir, scope+".enc")

			header := utils.MagicHeader(scope)
			body := []byte("SECRET-DATA")
			content := append(header, body...)

			err := os.WriteFile(path, content, 0644)
			assert.NoError(t, err)

			ok, err := utils.IsEncryptedFile(path, scope)
			assert.NoError(t, err)
			assert.True(t, ok)

			cf, err := utils.ReadEncryptedFile(path)
			assert.NoError(t, err)
			assert.Equal(t, body, cf.Data)
			assert.Equal(t, string(header), cf.Header)
		})
	}

	t.Run("invalid file", func(t *testing.T) {
		tmpDir := t.TempDir()
		path := filepath.Join(tmpDir, "bad.enc")

		err := os.WriteFile(path, []byte("BAD"), 0644)
		assert.NoError(t, err)

		_, err = utils.ReadEncryptedFile(path)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "encrypted signature")
	})
}
