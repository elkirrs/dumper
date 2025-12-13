package crypt_config

import (
	"context"
	"crypto/rand"
	"dumper/internal/domain/app"
	"dumper/pkg/utils"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/term"
)

type Crypt struct {
	ctx   context.Context
	flags *app.Flags
}

func NewApp(
	ctx context.Context,
	flags *app.Flags,
) *Crypt {
	return &Crypt{
		ctx:   ctx,
		flags: flags,
	}
}

func (c *Crypt) Run() error {
	switch c.flags.Mode {
	case "encrypt":
		if c.flags.Input == "" {
			c.flags.Input = "config.yaml"
		}

		valid := map[string]bool{
			"app":    true,
			"device": true,
			"both":   true,
		}

		if !valid[c.flags.Scope] {
			return fmt.Errorf("available Scope: app | device")
		}

		if c.flags.Password == "" {
			fmt.Println("Enter the password :")
			password, err := term.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				return fmt.Errorf("input error: %v", err)
			}
			c.flags.Password = strings.TrimSpace(string(password))
		}

		if err := c.encryptConfig(c.flags.Input, c.flags.Input, c.flags.Password); err != nil {
			fmt.Println("Encryption error:", err)
			return nil
		}
		fmt.Println("The configuration is encrypted in", c.flags.Input)

	case "decrypt":
		if c.flags.Input == "" {
			c.flags.Input = "config.yaml"
		}

		if c.flags.Password == "" {
			fmt.Println("Enter the password :")
			password, err := term.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				return fmt.Errorf("input error: %v", err)
			}
			c.flags.Password = strings.TrimSpace(string(password))
		}

		if err := c.decryptWithPassword(c.flags.Input, c.flags.Input, c.flags.Password); err != nil {
			fmt.Println("Error:", err)
			return nil
		}
		fmt.Println("The configuration is decrypted in", c.flags.Input)

	case "recovery":
		if c.flags.Input == "" {
			c.flags.Input = "config.yaml"
		}

		if c.flags.Recovery == "" {
			fmt.Println("Use: -crypt config -mode recovery -token <recovery key> -input config.yaml")
			return nil
		}

		if err := c.recoverConfig(c.flags.Input, c.flags.Input, c.flags.Recovery); err != nil {
			fmt.Println("Recovery error:", err)
			return nil
		}
		fmt.Println("The config was restored using the recovery token in", c.flags.Input)

	default:
		fmt.Println("Available modes: encrypt | decrypt | recovery")
	}

	return nil
}

func (c *Crypt) encryptConfig(input, output, password string) error {
	plain, err := os.ReadFile(input)
	if err != nil {
		return err
	}

	isEnc, err := utils.IsEncryptedFile(input, c.flags.Scope)
	if err == nil && isEnc {
		return fmt.Errorf("the %s file is already encrypted", input)
	}

	salt := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return err
	}

	secretCrypt := utils.SecretCrypt(c.flags)

	passwordKey := utils.DeriveKey(password, secretCrypt.SecretKey, salt)
	finalKey := utils.ComputeFinalKey(passwordKey, secretCrypt.DeviceKey)

	encConfig := utils.EncryptAES(finalKey, plain)

	deriveAppSecret := utils.DeriveAppKey(secretCrypt.SecretKey, secretCrypt.DeviceKey)
	encKeyForApp := utils.EncryptAES(deriveAppSecret, finalKey)

	data := append(utils.MagicHeader(c.flags.Scope), salt...)
	data = append(data, byte(len(encKeyForApp)))
	data = append(data, encKeyForApp...)
	data = append(data, encConfig...)

	recoveryToken := hex.EncodeToString(finalKey[:])
	fmt.Println("Recovery token (save securely, allows recovery on any device with password):", recoveryToken)

	return os.WriteFile(output, data, 0600)
}

func (c *Crypt) decryptWithPassword(input, output, password string) error {
	cFile, err := utils.ReadEncryptedFile(input)
	if err != nil {
		return err
	}

	data := cFile.Data
	salt := data[:16]
	keyLen := int(data[16])
	offset := 17 + keyLen
	encConfig := data[offset:]

	c.flags.Scope = utils.GetScope(cFile.Header)
	secretCrypt := utils.SecretCrypt(c.flags)
	passwordKey := utils.DeriveKey(password, secretCrypt.SecretKey, salt)
	finalKey := utils.ComputeFinalKey(passwordKey, secretCrypt.DeviceKey)

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
	cFile, err := utils.ReadEncryptedFile(input)
	if err != nil {
		return fmt.Errorf("failed to read encrypted file: %v", err)
	}

	data := cFile.Data

	if len(data) < 17 {
		return fmt.Errorf("invalid encrypted file (too short)")
	}

	keyLen := int(data[16])
	if len(data) < 17+keyLen {
		return fmt.Errorf("invalid encrypted file (encKey length mismatch)")
	}
	offset := 17 + keyLen
	encConfig := data[offset:]

	finalKey, err := hex.DecodeString(recoveryToken)
	if err != nil {
		return fmt.Errorf("invalid recovery token format: %v", err)
	}
	if len(finalKey) != 32 {
		return fmt.Errorf("invalid recovery token length: expected 32 bytes, got %d", len(finalKey))
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
