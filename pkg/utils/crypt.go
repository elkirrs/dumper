package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"os"

	"golang.org/x/crypto/pbkdf2"
)

func EncryptAES(key, data []byte) []byte {
	block, _ := aes.NewCipher(key)
	gcm, _ := cipher.NewGCM(block)
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err)
	}
	return append(nonce, gcm.Seal(nil, nonce, data, nil)...)
}

func DecryptAES(key, data []byte) ([]byte, error) {
	block, _ := aes.NewCipher(key)
	gcm, _ := cipher.NewGCM(block)
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

func DeriveKey(password string, appSecret, salt []byte) []byte {
	return pbkdf2.Key(append([]byte(password), appSecret...), salt, 100000, 32, sha256.New)
}

func DeriveAppKey(appSecret, deviceID []byte, salt []byte) []byte {

	// the ability to update file only in an environment where there was coding
	sum := sha256.Sum256(appSecret)
	return sum[:]

	// the ability to read and update file only in an environment where there was coding
	//combined := append(appSecret, deviceID...)
	//return pbkdf2.Key(combined, salt, 100000, 32, sha256.New)
}

func ComputeFinalKey(passwordKey, deviceKey []byte) []byte {
	return pbkdf2.Key(passwordKey, deviceKey, 100000, 32, sha256.New)
}

func IsEncrypted(data []byte) bool {
	header := MagicHeader()
	if len(data) < len(header) {
		return false
	}
	return string(data[:len(header)]) == string(header)
}

func LooksEncrypted(data []byte) bool {
	if len(data) < 32 {
		return false
	}
	for _, b := range data[:16] {
		if b < 32 || b > 126 {
			return true
		}
	}
	return false
}

func IsEncryptedFile(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	header := make([]byte, len(MagicHeader()))
	n, err := f.Read(header)
	if err != nil && err != io.EOF {
		return false, err
	}
	return n == len(MagicHeader()) && string(header) == string(MagicHeader()), nil
}

func ReadEncryptedFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if len(data) < len(MagicHeader()) || string(data[:len(MagicHeader())]) != string(MagicHeader()) {
		return nil, fmt.Errorf("the file does not have an encrypted signature")
	}
	return data[len(MagicHeader()):], nil
}

func MagicHeader() []byte {
	return []byte("ENCFDCA")
}
