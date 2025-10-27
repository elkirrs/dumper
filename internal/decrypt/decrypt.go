package decrypt

import (
	"bytes"
	"dumper/internal/command/encrypt"
	encOpts "dumper/internal/domain/encrypt"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Encrypt struct {
	FilePath string
	Password string
	Crypt    string
}

func NewApp(
	filePath string,
	password string,
	crypt string,
) *Encrypt {
	return &Encrypt{
		FilePath: filePath,
		Password: password,
		Crypt:    crypt,
	}
}

func (e *Encrypt) Decrypt() error {
	if e.FilePath == "" {
		return errors.New("-dec option is required")
	}

	if e.Password == "" && e.Crypt == "aes" {
		return errors.New("--pass option is required")
	}

	if e.Crypt == "" {
		return errors.New("-crypt option is required")
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

	return nil
}
