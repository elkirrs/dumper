package crypt_backup

import (
	"bytes"
	"dumper/internal/command/encrypt"
	encOpts "dumper/internal/domain/encrypt"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/term"
)

type CryptBackup struct {
	FilePath string
	Password string
	Crypt    string
	Type     string
}

func NewApp(
	filePath string,
	password string,
	crypt string,
	typeCrypt string,
) *CryptBackup {
	return &CryptBackup{
		FilePath: filePath,
		Password: password,
		Crypt:    crypt,
		Type:     typeCrypt,
	}
}

func (e *CryptBackup) Run() error {

	switch e.Type {
	case "decrypt":
		return e.Decrypt()
	case "encrypt":
		return e.Encrypt()
	default:
		return errors.New("unknown crypt type: " + e.Type)
	}

	return nil
}

func (e *CryptBackup) Decrypt() error {
	if e.FilePath == "" {
		return fmt.Errorf("file path is empty")
	}

	if e.Password == "" {
		fmt.Println("Enter the password :")
		password, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return fmt.Errorf("input error: %v", err)
		}

		e.Password = strings.TrimSpace(string(password))
	}

	encOption := encOpts.Options{
		Password: e.Password,
		FilePath: e.FilePath,
		Type:     e.Crypt,
		Crypt:    "decrypt",
	}

	encApp := encrypt.NewApp(&encOption)
	cmdData, err := encApp.Generate()
	if err != nil {
		return err
	}

	args := strings.Fields(cmdData.CMD)
	bin := args[0]
	args = args[1:]

	binPath, err := exec.LookPath(bin)
	if err != nil {
		return fmt.Errorf("cannot find binary %s: %v", bin, err)
	}

	cmd := exec.Command(binPath, args...)
	cmd.Env = os.Environ()

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("decryption failed: %v\nOutput: %s", err, out.String())
	}

	fmt.Println("Decryption succeeded")
	return nil
}

func (e *CryptBackup) Encrypt() error {
	if e.FilePath == "" {
		return fmt.Errorf("file path is empty")
	}

	if e.Password == "" {
		fmt.Println("Enter the password :")
		password, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return fmt.Errorf("input error: %v", err)
		}

		e.Password = strings.TrimSpace(string(password))
	}

	encOption := encOpts.Options{
		Password: e.Password,
		FilePath: e.FilePath,
		Type:     e.Crypt,
		Crypt:    "encrypt",
	}

	encApp := encrypt.NewApp(&encOption)
	cmdData, err := encApp.Generate()
	if err != nil {
		return err
	}

	args := strings.Fields(cmdData.CMD)
	bin := args[0]
	args = args[1:]

	binPath, err := exec.LookPath(bin)
	if err != nil {
		return fmt.Errorf("cannot find binary %s: %v", bin, err)
	}

	cmd := exec.Command(binPath, args...)
	cmd.Env = os.Environ()

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("encryption failed: %v\nOutput: %s", err, out.String())
	}

	fmt.Println("Encryption succeeded")
	return nil
}
