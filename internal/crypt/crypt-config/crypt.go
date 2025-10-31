package crypt_config

import (
	"context"
	"crypto/rand"
	"dumper/pkg/utils"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

type Crypt struct {
	ctx         context.Context
	mode        string
	input       string
	password    string
	appSecret   string
	recoveryKey string
}

func NewApp(
	ctx context.Context,
	mode string,
	input string,
	password string,
	appSecret string,
	recoveryKey string,
) *Crypt {
	return &Crypt{
		ctx:         ctx,
		input:       input,
		mode:        mode,
		password:    password,
		appSecret:   appSecret,
		recoveryKey: recoveryKey,
	}
}

func (c *Crypt) Run() error {
	switch c.mode {
	case "encrypt":
		if c.password == "" || c.input == "" {
			fmt.Println("Use: -crypt config -mode encrypt -password <password> -input config.yaml")
			return nil
		}
		if err := c.encryptConfig(c.input, c.input, c.password); err != nil {
			fmt.Println("Encryption error:", err)
			return nil
		}
		fmt.Println("The configuration is encrypted in", c.input)

	case "decrypt":
		if c.password == "" || c.input == "" {
			fmt.Println("Use: -crypt config -mode decrypt -password <password> -input config.yaml")
			return nil
		}
		if err := c.decryptWithPassword(c.input, c.input, c.password); err != nil {
			fmt.Println("Error:", err)
			return nil
		}
		fmt.Println("The configuration is decrypted in", c.input)

	case "recover":
		if c.input == "" || c.recoveryKey == "" {
			fmt.Println("Use: -crypt config -mode recover -recovery <recovery key> -input config.yaml")
			return nil
		}

		if err := c.recoverConfig(c.input, c.input, c.recoveryKey); err != nil {
			fmt.Println("Recovery error:", err)
			return nil
		}
		fmt.Println("The config was restored using the recovery token in", c.input)

	default:
		fmt.Println("Available modes: encrypt | decrypt | recover")
	}

	return nil
}

func (c *Crypt) encryptConfig(input, output, password string) error {
	plain, err := os.ReadFile(input)
	if err != nil {
		return err
	}

	isEnc, err := utils.IsEncryptedFile(input)
	if err == nil && isEnc {
		return fmt.Errorf("the %s file is already encrypted", input)
	}

	salt := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return err
	}

	appSecret := []byte(c.appSecret)
	deviceKey := utils.GetDeviceKey()
	passwordKey := utils.DeriveKey(password, appSecret, salt)
	finalKey := utils.ComputeFinalKey(passwordKey, deviceKey)

	encConfig := utils.EncryptAES(finalKey, plain)

	deriveAppSecret := utils.DeriveAppKey(appSecret, deviceKey, salt)
	encKeyForApp := utils.EncryptAES(deriveAppSecret, finalKey)

	data := append(utils.MagicHeader(), salt...)
	data = append(data, byte(len(encKeyForApp)))
	data = append(data, encKeyForApp...)
	data = append(data, encConfig...)

	recoveryToken := hex.EncodeToString(finalKey[:])
	fmt.Println("Recovery token (save securely, allows recovery on any device with password):", recoveryToken)

	return os.WriteFile(output, data, 0600)
}

func (c *Crypt) decryptWithPassword(input, output, password string) error {
	data, err := utils.ReadEncryptedFile(input)
	if err != nil {
		return err
	}

	salt := data[:16]
	keyLen := int(data[16])
	offset := 17 + keyLen
	encConfig := data[offset:]

	passwordKey := utils.DeriveKey(password, []byte(c.appSecret), salt)
	deviceKey := utils.GetDeviceKey()
	finalKey := utils.ComputeFinalKey(passwordKey, deviceKey)

	plain, err := utils.DecryptAES(finalKey, encConfig)
	if err != nil {
		return fmt.Errorf("wrong password or wrong device: %v", err)
	}

	return os.WriteFile(output, plain, 0600)
}

func (c *Crypt) recoverConfig(input, output, recoveryToken string) error {
	if input == "" || recoveryToken == "" {
		return fmt.Errorf("input and recovery token must be specified")
	}

	data, err := os.ReadFile(input)
	if err != nil {
		return fmt.Errorf("failed to read input file: %v", err)
	}

	keyLen := int(data[16])
	offset := 17 + keyLen
	encConfig := data[offset:]

	finalKey, err := hex.DecodeString(recoveryToken)
	if err != nil {
		return fmt.Errorf("invalid recovery token format: %v", err)
	}

	plain, err := utils.DecryptAES(finalKey, encConfig)
	if err != nil {
		return fmt.Errorf("decryption failed with recovery token: %v", err)
	}

	if err := os.WriteFile(output, plain, 0600); err != nil {
		return fmt.Errorf("failed to write output: %v", err)
	}

	return nil
}
