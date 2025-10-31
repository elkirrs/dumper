package utils_test

import (
	"dumper/pkg/utils"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncryptDecryptAES(t *testing.T) {
	key := []byte("12345678901234567890123456789012") // 32 byte for AES-256
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
	salt := []byte("salt")

	key := utils.DeriveAppKey(appSecret, deviceID, salt)
	assert.Len(t, key, 32)
}

func TestComputeFinalKey_Consistency(t *testing.T) {
	passwordKey := []byte("passwordkey1234567890123456")
	deviceKey := []byte("devicekey1234567890123456789")

	key1 := utils.ComputeFinalKey(passwordKey, deviceKey)
	key2 := utils.ComputeFinalKey(passwordKey, deviceKey)
	assert.Equal(t, key1, key2)
}

func TestIsEncrypted(t *testing.T) {
	header := utils.MagicHeader()
	data := append(header, []byte("payload")...)
	assert.True(t, utils.IsEncrypted(data), "data with magic header should be encrypted")

	data2 := []byte("abcdpayload")
	assert.False(t, utils.IsEncrypted(data2), "data without magic header should not be encrypted")
}

func TestLooksEncrypted(t *testing.T) {
	data := make([]byte, 32)
	for i := 0; i < 16; i++ {
		data[i] = 200 // not ascii
	}
	assert.True(t, utils.LooksEncrypted(data))

	data2 := []byte("this is normal ascii data................")
	assert.False(t, utils.LooksEncrypted(data2))
}

func TestIsEncryptedFileAndReadEncryptedFile(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "file.enc")

	content := append(utils.MagicHeader(), []byte("secretdata")...)
	err := os.WriteFile(path, content, 0644)
	assert.NoError(t, err)

	ok, err := utils.IsEncryptedFile(path)
	assert.NoError(t, err)
	assert.True(t, ok)

	data, err := utils.ReadEncryptedFile(path)
	assert.NoError(t, err)
	assert.Equal(t, []byte("secretdata"), data)

	path2 := filepath.Join(tmpDir, "file2.enc")
	_ = os.WriteFile(path2, []byte("noheader"), 0644)
	_, err = utils.ReadEncryptedFile(path2)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "encrypted signature")
}
