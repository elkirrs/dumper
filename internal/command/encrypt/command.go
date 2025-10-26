package encrypt

import (
	"dumper/internal/command/encrypt/aes"
	"dumper/internal/domain/encrypt"
	"fmt"
)

type Encrypt struct {
	Options *encrypt.Options
}

type EncryptGenerator interface {
	Generate(opts *encrypt.Options) *encrypt.DataCrypt
}

func NewApp(options *encrypt.Options) *Encrypt {
	return &Encrypt{
		Options: options,
	}
}

func (e *Encrypt) Generate() (*encrypt.DataCrypt, error) {

	var gen EncryptGenerator

	switch e.Options.Type {
	case "aes":
		gen = aes.AESGenerator{}
	default:
		return nil, fmt.Errorf("unknown encrypt type: %s", e.Options.Type)
	}

	cmdData := gen.Generate(e.Options)

	return cmdData, nil
}
