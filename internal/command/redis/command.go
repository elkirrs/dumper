package redis

import (
	"dumper/internal/command"
	cmdCfg "dumper/internal/domain/command-config"
	"dumper/internal/domain/config/setting"
	"fmt"
)

type RedisGenerator struct{}

func (g RedisGenerator) Generate(data *cmdCfg.ConfigData, settings *setting.Settings) (string, string) {
	ext := "rdb"

	fileName := fmt.Sprintf("%s.%s", data.DumpName, ext)
	remotePath := fmt.Sprintf("./%s", fileName)

	host := "127.0.0.1"

	baseCmd := fmt.Sprintf(
		"redis-cli -h %s -p %s -a %s --rdb %s",
		host, data.Port, data.Password, remotePath,
	)

	if data.Options.Mode == "save" {
		saveCmd := fmt.Sprintf(
			"redis-cli -h %s -p %s -a %s SAVE && ",
			host, data.Port, data.Password,
		)
		baseCmd = saveCmd + baseCmd
	}

	if *settings.Archive {
		baseCmd = fmt.Sprintf("%s | gzip", baseCmd)
		ext += ".gz"
		remotePath += ".gz"
	}

	if settings.DumpLocation == "server" {
		return baseCmd, remotePath
	}

	return baseCmd, remotePath
}

func init() {
	command.Register("redis", RedisGenerator{})
}
