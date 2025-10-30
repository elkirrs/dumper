package crypt

import (
	"context"
	cryptBackup "dumper/internal/crypt/crypt-backup"
	cryptConfig "dumper/internal/crypt/crypt-config"
	"dumper/internal/domain/app"
	"dumper/pkg/utils"
	"fmt"
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

	switch c.flags.Crypt {
	case "backup":
		cryptBackupApp := cryptBackup.NewApp(
			c.flags.Input,
			c.flags.Password,
			"aes",
		)
		err := cryptBackupApp.Decrypt()
		if err != nil {
			return err
		}

	case "config":
		cryptConfigApp := cryptConfig.NewApp(
			c.ctx,
			c.flags.Mode,
			c.flags.Input,
			c.flags.Password,
			c.flags.AppSecret,
			c.flags.Recovery,
		)
		err := cryptConfigApp.Run()
		if err != nil {
			return err
		}

	default:
		return fmt.Errorf("unknown crypt flag: %s", c.flags.Crypt)
	}

	return nil
}

func DecryptInApp(data []byte, appSecret string) ([]byte, error) {
	keyLen := int(data[16])
	offset := 17
	encKey := data[offset : offset+keyLen]
	offset += keyLen
	encConfig := data[offset:]

	salt := data[:16]
	deviceKey := utils.GetDeviceKey()
	deriveAppSecret := utils.DeriveAppKey([]byte(appSecret), deviceKey, salt)

	finalKey, err := utils.DecryptAES(deriveAppSecret, encKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt key: %v", err)
	}

	plain, err := utils.DecryptAES(finalKey, encConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt config: %v", err)
	}

	return plain, nil
}
