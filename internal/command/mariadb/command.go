package mariadb

import (
	"dumper/internal/command"
	cmdCfg "dumper/internal/domain/command-config"
	"dumper/internal/domain/config/setting"
	"fmt"
)

type MariaDBGenerator struct{}

func (g MariaDBGenerator) Generate(data *cmdCfg.ConfigData, settings *setting.Settings) (string, string) {
	if data.Port == "" {
		data.Port = "3306"
	}

	ext := "sql"

	baseCmd := fmt.Sprintf(
		"/usr/bin/mariadb-dump -u%s -p%s -h127.0.0.1 -P%s %s",
		data.User, data.Password, data.Port, data.Name,
	)

	if *settings.Archive {
		baseCmd += " | gzip"
		ext += ".gz"
	}

	fileName := fmt.Sprintf("%s.%s", data.DumpName, ext)
	remotePath := fmt.Sprintf("./%s", fileName)

	if settings.DumpLocation == "server" {
		return fmt.Sprintf("%s > %s", baseCmd, remotePath), remotePath
	}

	return baseCmd, remotePath
}

func init() {
	command.Register("mariadb", MariaDBGenerator{})
}
