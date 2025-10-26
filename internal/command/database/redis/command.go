package redis

import (
	command "dumper/internal/command/database"
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
		"redis-cli -h %s -p %s -a %s --rdb",
		host, data.Port, data.Password,
	)

	if data.Options.Mode == "save" {
		saveCmd := fmt.Sprintf(
			"redis-cli -h %s -p %s -a %s SAVE && ",
			host, data.Port, data.Password,
		)
		baseCmd = saveCmd + baseCmd
	}

	if *settings.Archive {
		remotePath += ".gz"
		baseCmd = fmt.Sprintf("%s - | gzip > %s", baseCmd, remotePath)
	} else {
		baseCmd = fmt.Sprintf("%s %s", baseCmd, remotePath)
	}

	if settings.DumpLocation == "server" {
		return baseCmd, remotePath
	}

	return baseCmd, remotePath
}

func init() {
	command.Register("redis", RedisGenerator{})
}
