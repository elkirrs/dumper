package aes

import (
	"dumper/internal/domain/encrypt"
	"dumper/pkg/utils/suffix"
	"fmt"
)

type AESGenerator struct{}

func (g AESGenerator) Generate(opts *encrypt.Options) *encrypt.DataCrypt {

	out := &encrypt.DataCrypt{}

	switch opts.Crypt {
	case "encrypt":
		encPath := opts.FilePath + ".enc"
		out.CMD = encryptFile(opts.FilePath, encPath, opts.Password)
		out.Name = encPath

	case "decrypt":
		decPath := suffix.RemoveSuffix(opts.FilePath, ".enc")
		out.CMD = decryptFile(opts.FilePath, decPath, opts.Password)
		out.Name = decPath
	}

	return out
}

func encryptFile(remotePath, encPath, password string) string {
	return fmt.Sprintf(
		"openssl enc -aes-256-cbc -salt -pbkdf2 -iter 100000 -in %s -out %s -k %s",
		remotePath,
		encPath,
		password,
	)
}

func decryptFile(remotePath, decPath, password string) string {
	return fmt.Sprintf(
		"openssl enc -d -aes-256-cbc -pbkdf2 -iter 100000 -in %s -out %s -k %s",
		remotePath,
		decPath,
		password,
	)
}
