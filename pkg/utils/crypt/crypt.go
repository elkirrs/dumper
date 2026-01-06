package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	app2 "dumper/internal/domain/app"
	device2 "dumper/pkg/utils/device"
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

func DeriveAppKey(secretKey, deviceKey []byte) []byte {
	combined := append(secretKey, deviceKey...)
	sum := sha256.Sum256(combined)
	return sum[:]
}

func ComputeFinalKey(passwordKey, deviceKey []byte) []byte {
	return pbkdf2.Key(passwordKey, deviceKey, 100000, 32, sha256.New)
}

func IsEncrypted(data []byte) bool {
	mHeader := data[:magicHeaderLength]
	scope := GetScope(string(mHeader))
	header := MagicHeader(scope)

	if len(header) != magicHeaderLength {
		return false
	}

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

func IsEncryptedFile(path, scope string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	header := make([]byte, len(MagicHeader(scope)))

	n, err := f.Read(header)
	if err != nil && err != io.EOF {
		return false, err
	}
	return len(header) != 0 && n == len(MagicHeader(scope)) && string(header) == string(MagicHeader(scope)), nil
}

type CFile struct {
	Data   []byte
	Header string
}

func ReadEncryptedFile(path string) (CFile, error) {
	data, err := os.ReadFile(path)
	mHeader := string(data[:magicHeaderLength])
	scope := GetScope(mHeader)

	var cFile CFile
	if err != nil {
		return cFile, err
	}
	if len(data) < len(MagicHeader(scope)) || string(data[:len(MagicHeader(scope))]) != string(MagicHeader(scope)) {
		return cFile, fmt.Errorf("the file does not have an encrypted signature")
	}
	cFile.Data = data[magicHeaderLength:]
	cFile.Header = mHeader

	return cFile, nil
}

var (
	magicHeaderLength = 7
	both              = "ENCFDCA"
	app               = "ENCADCA"
	device            = "ENCDDCA"
)

func MagicHeader(scope string) []byte {
	magicHeader := ""
	switch scope {
	case "app":
		magicHeader = app
	case "device":
		magicHeader = device
	default:
		magicHeader = both
	}
	return []byte(magicHeader)
}

func GetScope(magicHeader string) string {
	switch magicHeader {
	case app:
		return "app"
	case device:
		return "device"
	case both:
		return "both"
	default:
		return ""
	}
}

type SCrypt struct {
	SecretKey []byte
	DeviceKey []byte
}

func SecretCrypt(flag *app2.Flags) SCrypt {
	var appSecret []byte
	var deviceKey []byte

	switch flag.Scope {
	case "app":
		appSecret = []byte(flag.AppSecret)
	case "device":
		deviceKey = device2.GetDeviceKey()
	default:
		appSecret = []byte(flag.AppSecret)
		deviceKey = device2.GetDeviceKey()
	}

	return SCrypt{
		SecretKey: appSecret,
		DeviceKey: deviceKey,
	}
}
